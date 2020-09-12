package uniq

import "testing"

var sampleText = []string{
	"I love music.\n",
	"I love music.\n",
	"I love music.\n",
	"\n",
	"I love music of Kartik.\n",
	"I love music of Kartik.\n",
	"Thanks.",
}

var caseInsesitiveText = []string{
	"I LOVE MUSIC.\n",
	"I love music.\n",
	"I LoVe MuSiC.\n",
	"\n",
	"I love MuSIC of Kartik.\n",
	"I love music of Kartik.\n",
	"Thanks.",
}



func TestUniqSuccess(t *testing.T) {
	var TestCases = []struct{
		in []string
		flags OptionParams
		out []string
	}{
		{sampleText,
			OptionParams{},
		[]string{
			"I love music.\n",
			"\n",
			"I love music of Kartik.\n",
			"Thanks.",
			},
		},
		{
			sampleText,
			OptionParams{C: true},
			[]string{
				"3 I love music.\n",
				"1 \n",
				"2 I love music of Kartik.\n",
				"1 Thanks.",
			},
		},
		{
			sampleText,
			OptionParams{D: true},
			[]string{
				"I love music.\n",
				"I love music of Kartik.\n",
			},
		},
		{
			sampleText,
			OptionParams{U: true},
			[]string{
				"\n",
				"Thanks.",
			},
		},
		{caseInsesitiveText,
			OptionParams{I: true},
			[]string{
				"I LOVE MUSIC.\n",
				"\n",
				"I love MuSIC of Kartik.\n",
				"Thanks.",
			},
		},
		{[]string {
				"We love music.\n",
				"I love music.\n",
				"They love music.\n",
				"\n",
				"I love music of Kartik.\n",
				"We love music of Kartik.\n",
				"Thanks.",
			},
			OptionParams{F: 1},
			[]string{
				"We love music.\n",
				"\n",
				"I love music of Kartik.\n",
				"Thanks.",
			},
		},
		{[]string {
			"I love music.\n",
			"A love music.\n",
			"C love music.\n",
			"\n",
			"I love music of Kartik.\n",
			"We love music of Kartik.\n",
			"Thanks.",
			},
			OptionParams{S: 1},
			[]string{
				"I love music.\n",
				"\n",
				"I love music of Kartik.\n",
				"We love music of Kartik.\n",
				"Thanks.",
			},
		},
		{[]string {
			"We are not love music.\n",
			"They are all love music.\n",
			"What does the love music.\n",
			"\n",
			"I love sth of Kartik.\n",
			"We love sth of Kartik.\n",
			"Thanks.",
		},
			OptionParams{F: 2, S: 3, C: true},
			[]string{
				"3 We are not love music.\n",
				"1 \n",
				"2 I love sth of Kartik.\n",
				"1 Thanks.",
			},
		},
		{
			caseInsesitiveText,
			OptionParams{I: true, F: 2, D: true},
			[]string{
				"I LOVE MUSIC.\n",
				"I love MuSIC of Kartik.\n",
			},
		},
		{[]string {
			"We    love music.\n",
			"I     love     music.\n",
			"They love     music.\n",
			"\n",
			"I love music of Kartik.\n",
			"We love music of Kartik.\n",
			"Thanks.",
			},
			OptionParams{F: 1},
			[]string{
				"We    love music.\n",
				"\n",
				"I love music of Kartik.\n",
				"Thanks.",
			},
		},

	}

	for _, tCase := range TestCases {
		t.Run("", func(t *testing.T) {
			res, err := UniqUtility(tCase.in, &tCase.flags)
			if err != nil {
				t.Errorf("got err %e\n", err)
			}
			if len(res) != len(tCase.out) {
				t.Errorf("got len %d, want len %d", len(res), len(tCase.out))
			}
			for i, str := range res {
				if str != tCase.out[i] {
					t.Errorf("got %s, want %s", str, tCase.out[i])
				}
			}
		})
	}
}

func TestUniqFail(t *testing.T) {
	var TestCases = []struct{
		in []string
		flags OptionParams
		out []string
	}{
		{sampleText,
			OptionParams{C: true, D: true, U: true},
			[]string{},
		},
	}

	for _, tCase := range TestCases {
		t.Run("", func(t *testing.T) {
			res, err := UniqUtility(tCase.in, &tCase.flags)
			if err == nil {
				t.Errorf("want err")
			}
			if len(res) != len(tCase.out) {
				t.Errorf("got len %d, want len %d", len(res), len(tCase.out))
			}
			for i, str := range res {
				if str != tCase.out[i] {
					t.Errorf("got %s, want %s", str, tCase.out[i])
				}
			}
		})
	}
}
