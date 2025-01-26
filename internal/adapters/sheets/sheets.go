package sheets

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/guregu/null/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/internal/models"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type SheetClient struct {
	service         *sheets.Service
	credentialsPath string

	log log.Logger
}

func New(logger log.Logger, credentialsPath string) *SheetClient {
	return &SheetClient{credentialsPath: credentialsPath, log: logger}
}

func (s *SheetClient) Start(ctx context.Context) error {
	b, err := os.ReadFile(s.credentialsPath + "credentials.json")
	if err != nil {
		return errors.Wrap(err, "Unable to read client secret file")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return errors.Wrap(err, "Unable to parse client secret file to config")
	}

	client, err := getClient(ctx, s.credentialsPath, config)
	if err != nil {
		return errors.Wrap(err, "Unable to get client")
	}

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	s.service = service

	return nil
}

func (s *SheetClient) Stop() error {
	return nil
}

func (s *SheetClient) HasSheet(id string) (bool, error) {
	sheet, err := s.service.Spreadsheets.Get(id).Do()
	if err != nil {
		return false, err
	}

	return sheet != nil, nil
}

func (s *SheetClient) GetPools(ctx context.Context, sheetID string) ([]models.Pool, error) {
	result, err := s.service.Spreadsheets.Get(sheetID).Context(ctx).IncludeGridData(true).Do()
	if err != nil {
		return nil, err
	}

	pools := make([]models.Pool, 0, len(result.Sheets))

	for _, sh := range result.Sheets {
		if sh.Properties.Hidden {
			continue
		}

		if len(sh.Data) < 1 || sh.Data[0] == nil {
			continue
		}

		data := sh.Data[0]

		if len(data.RowData) < 2 {
			continue
		}

		rowData := data.RowData[1]

		if len(rowData.Values) < 1 || rowData.Values[0] == nil {
			continue
		}

		value, err := strconv.ParseFloat(rowData.Values[0].FormattedValue, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert value to int")
		}

		measurements := make([]common.Measurement, 0, len(data.RowData)-1)
		actions := make([]common.Action, 0, len(data.RowData)-1)
		chemicals := make([]common.Chemicals, 0, len(data.RowData)-1)

		year := 2024

		for i, v := range data.RowData {
			if i == 0 {
				continue
			}
			lg := s.log.WithField("row", i)

			if len(v.Values) < 2 || v.Values[1] == nil {
				lg.Info("skipping row")
				continue
			}

			date, err := parseDate(v.Values[1].FormattedValue, year)
			if err != nil {
				lg.WithError(err).Error("failed to parse date")
				continue
			}

			if date.Sub(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) > 0 {
				year = 2025
			}

			measurement := getMeasurement(v, lg)
			if measurement != nil {
				measurement.CreatedAt = date
				measurements = append(measurements, *measurement)
			}

			action := getAction(v)
			if action != nil {
				action.CreatedAt = date
				actions = append(actions, *action)
			}

			chemical := getChemicals(v, lg)
			if chemical != nil {
				chemical.CreatedAt = date
				chemicals = append(chemicals, *chemical)
			}

		}

		pools = append(pools, models.Pool{
			Pool: common.PoolData{
				Name:   sh.Properties.Title,
				Volume: value,
			},
			Measurements: measurements,
			Chemicals:    chemicals,
			Actions:      actions,
		})
	}

	return pools, nil
}

func getAction(row *sheets.RowData) *common.Action {
	sheetIndexToAction := map[int]common.ActionType{
		23: common.ActionVacuum,
		24: common.ActionNet,
		25: common.ActionBrush,
		26: common.ActionBackwash,
		27: common.ActionSkimmerBasketClean,
		28: common.ActionPumpBasketClean,
		29: common.ActionScumLine,
	}

	actionTypes := make([]common.ActionType, 0, len(sheetIndexToAction))

	for i, action := range sheetIndexToAction {
		if row.Values[i] != nil && row.Values[i].FormattedValue != "" {
			actionTypes = append(actionTypes, action)
		}
	}

	if len(actionTypes) > 0 {
		return &common.Action{
			Types: actionTypes,
		}
	}

	return nil
}

func getChemicals(row *sheets.RowData, lg log.Logger) *common.Chemicals {
	chemicals := make(map[common.ChemicalProduct]float64)

	for i := 9; i < 18; i++ {
		product, value := getChemicalProduct(row.Values, i, lg)

		if value.Valid {
			chemicals[product] = value.Float64
		}
	}
	i := 22

	product, value := getChemicalProduct(row.Values, i, lg)

	if value.Valid {
		chemicals[product] = value.Float64
	}

	if len(chemicals) > 0 {
		return &common.Chemicals{
			Products: chemicals,
		}
	}

	return nil
}

func getChemicalProduct(values []*sheets.CellData, i int, lg log.Logger) (common.ChemicalProduct, null.Float) {
	if values[i] == nil {
		return 0, null.Float{}
	}

	sheetIndexToChemicalProduct := map[int]common.ChemicalProduct{
		9:  common.CalciumHypochlorite65Percent,
		10: common.SodiumHypochlorite12Percent,
		11: common.SodiumHypochlorite14Percent,
		12: common.TCCA90PercentTablets,
		13: common.MultiActionTablets,
		14: common.TCCA90PercentGranules,
		15: common.Dichlor65Percent,
		16: common.HydrochloricAcid,
		17: common.SodiumBisulphate,
		22: common.SodiumBicarbonate,
	}

	product, ok := sheetIndexToChemicalProduct[i]
	if !ok {
		lg.WithField("product", values[i].FormattedValue).Warn("unknown chemical product")

		return 0, null.Float{}
	}

	return product, getFloatFromCells(values, i, lg)
}

func getMeasurement(row *sheets.RowData, lg log.Logger) *common.Measurement {
	chlorine := getFloatFromCells(row.Values, 2, lg)
	ph := getFloatFromCells(row.Values, 4, lg)
	alkalinity := getFloatFromCells(row.Values, 6, lg)

	if chlorine.Valid || ph.Valid || alkalinity.Valid {
		return &common.Measurement{
			Chlorine:   chlorine,
			PH:         ph,
			Alkalinity: alkalinity,
		}
	}

	return nil
}

func getFloatFromCells(cells []*sheets.CellData, i int, lg log.Logger) null.Float {
	value, err := getFloat(cells[i].FormattedValue)
	if err != nil {
		lg.WithError(err).WithField("column", i).Error("failed to convert value to float")
	}

	return value
}

func getFloat(value string) (null.Float, error) {
	if value == "" {
		return null.Float{}, nil
	}

	val, err := strconv.ParseFloat(strings.Replace(value, ",", ".", 1), 64)
	if err != nil {
		return null.Float{}, errors.Wrap(err, "failed to convert value to float")
	}

	return null.FloatFrom(val), nil
}

func parseDate(value string, year int) (time.Time, error) {
	date, err := time.Parse("02/01", value)
	if err == nil {
		date = date.AddDate(year, 0, 0)
		return date, nil
	}

	date, err = time.Parse("2/01", value)
	if err == nil {
		date = date.AddDate(year, 0, 0)
		return date, nil
	}

	date, err = time.Parse("02.01", value)
	if err == nil {
		date = date.AddDate(year, 0, 0)
		return date, nil
	}

	date, err = time.Parse("02/01/2006", value)
	if err == nil {
		return date, nil
	}

	return date, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, path string, config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := path + "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
