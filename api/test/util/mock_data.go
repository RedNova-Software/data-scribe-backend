package util_test

import "api/shared/models"

// Mock data for testing
func mockStaticData() *models.ReportSection {
	reportSection := models.ReportSection{
		Title:           "Moving Section",
		OutputGenerated: true,
		Questions: []models.ReportQuestion{
			{
				Label:    "q1",
				Question: "What is blue?",
				Answer:   "Blue is my favorite color!",
			},
			{
				Label:    "q2",
				Question: "What is red?",
				Answer:   "Red is also my favorite color!",
			},
			{
				Label:    "q3",
				Question: "What is green?",
				Answer:   "Green, green is the best.",
			},
		},
		CSVData: []models.ReportCSVData{
			{
				Label:           "csv1",
				Description:     "Average response time for station 1,2,3 (average for all 3)",
				OperationType:   "Average",
				OperationColumn: "Travel Time",
				AcceptedValues:  nil,
				FilterColumns: map[string][]string{
					"Station": {"1", "2", "3"},
					"Year":    {"2016", "2017", "2018"},
				},
				Result: "218.477",
			},
			{
				Label:           "csv2",
				Description:     "Number of fire incidents in 2016 for station 1 and district 4",
				OperationType:   "SetElementOccurences",
				OperationColumn: "Incident Type",
				AcceptedValues:  []string{"Fire"},
				FilterColumns: map[string][]string{
					"District":      {"4"},
					"Incident Type": {"Fire"},
					"Station":       {"1"},
					"Year":          {"2016"},
					"Address":       {"281 WOODLAWN RD"},
				},
				Result: "1",
			},
			{
				Label:           "csv3",
				Description:     "Number of vital signs absent and alcohol or drug-related incidents",
				OperationType:   "SetElementOccurences",
				OperationColumn: "Incident Type",
				AcceptedValues:  []string{"Vital signs absent, DOA", "Alcohol or drug related"},
				FilterColumns:   nil,
				Result:          "3474",
			},
		},
		TextOutputs: []models.ReportTextOutput{
			{
				Title:  "Text 1",
				Type:   "Static",
				Input:  "This is test text 1, we're going to splice in every question and csv data label\nq1: **q1\nq2: **q2\nq3: **q3\n\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3",
				Result: "",
			},
			{
				Title:  "Text 2",
				Type:   "Static",
				Input:  "This is test text 2, we're going to splice in every question and csv data label\n\ncsv1: **csv1\ncsv2: **csv2\ncsv3: **csv3\n\nq1: **q1\nq2: **q2\nq3: **q3\n\n\n",
				Result: "",
			},
		},
		ChartOutputs: []models.ReportChartOutput{
			{
				Title:                  "Travel & Turnout Time by Station",
				Type:                   "Bar",
				Description:            "A chart to show the travel time & Turnout time by station",
				XAxisTitle:             "Station",
				YAxisTitle:             "Time (s)",
				CartesianGrid:          true,
				IndependentColumnLabel: "Station",
				IndependentColumn:      "Station",
				AcceptedValues:         nil,
				FilterColumns:          nil,
				DependentColumns: []models.ReportOneDimConfig{
					{
						AggregateValueLabel: "Travel Time",
						Column:              "Travel Time",
						Description:         "This will get the average travel time by station",
						OperationType:       "Average",
						AcceptedValues:      nil,
						FilterColumns:       nil,
					},
					{
						AggregateValueLabel: "Turnout Time",
						Column:              "Turnout Time",
						Description:         "This will get the average turnout time by station",
						OperationType:       "Average",
						AcceptedValues:      nil,
						FilterColumns:       nil,
					},
				},
				Results: nil,
			},
			{
				Title:                  "Incidence per Station Per Year",
				Type:                   "Area",
				Description:            "A stacked area chart that will graph the number of incidents per station per year",
				XAxisTitle:             "Year",
				YAxisTitle:             "Number of Incidents",
				CartesianGrid:          true,
				IndependentColumnLabel: "Year",
				IndependentColumn:      "Year",
				AcceptedValues:         []string{"2016", "2017", "2018", "2019"},
				FilterColumns:          nil,
				DependentColumns: []models.ReportOneDimConfig{
					{
						AggregateValueLabel: "Station 1",
						Column:              "Incident Type",
						Description:         "Area for station 1 incidents over the years",
						OperationType:       "SetElementOccurences",
						AcceptedValues:      nil,
						FilterColumns: map[string][]string{
							"Station": {"1"},
						},
					},
					{
						AggregateValueLabel: "Station 2",
						Column:              "Incident Type",
						Description:         "Area for station 1 incidents over the years",
						OperationType:       "SetElementOccurences",
						AcceptedValues:      nil,
						FilterColumns: map[string][]string{
							"Station": {"2"},
						},
					},
				},
				Results: []map[string]interface{}{
					{
						"Station 1": 1486,
						"Station 2": 1599,
						"Year":      "2016",
					},
					{
						"Station 1": 1683,
						"Station 2": 1819,
						"Year":      "2017",
					},
					{
						"Station 1": 1838,
						"Station 2": 1930,
						"Year":      "2018",
					},
					{
						"Station 1": 1840,
						"Station 2": 1965,
						"Year":      "2019",
					},
				},
			},
		},
	}

	return &reportSection
}

func mockGeneratorData() *models.ReportSection {
	section := &models.ReportSection{
		Title: "Section Two - Generator",
		Questions: []models.ReportQuestion{
			{
				Label:    "questionOne",
				Question: "What's your favourite color?",
				Answer:   "",
			},
			{
				Label:    "questionTwo",
				Question: "What's your favourite city?",
				Answer:   "",
			},
		},
		TextOutputs: []models.ReportTextOutput{
			{
				Title:  "Generator One",
				Input:  "Tell me about this color: **questionOne",
				Type:   models.Generator,
				Result: "",
			},
			{
				Title:  "Generator Two",
				Input:  "Tell me about this city: **questionTwo",
				Type:   models.Generator,
				Result: "",
			},
		},
	}

	return section
}
