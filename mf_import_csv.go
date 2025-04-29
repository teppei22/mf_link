package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	// アカウント情報とインポート先URL
	url := "https://moneyforward.com/accounts/show_manual/xxxxxxxxxxxxxxx" // インポート先の口座URL
	user := "<自分のアカウント>"
	password := "<自分のパスワード>"

	// コマンドライン引数チェック
	if len(os.Args) != 2 {
		fmt.Println("No input_file!")
		fmt.Println("usage: go run mf_import_csv.go data_file.csv")
		os.Exit(1)
	}
	inputFile := os.Args[1]

	fmt.Println("Start :" + inputFile)

	// Selenium WebDriverの設定
	const (
		// ChromeDriverサービスのポート
		port = 4444
	)

	// Chromeのオプション設定
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// Chromeの追加設定
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
			"--disable-dev-shm-usage",
		},
	}
	caps.AddChrome(chromeCaps)

	// WebDriverサービスの開始
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
		os.Exit(1)
	}
	defer service.Stop()

	// WebDriverの作成
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Printf("Error connecting to the WebDriver: %v", err)
		os.Exit(1)
	}
	defer wd.Quit()

	// ブラウザの暗黙的待機時間設定
	wd.SetImplicitWaitTimeout(10 * time.Second)

	// マネーフォワードへのログイン
	err = wd.Get(url)
	if err != nil {
		log.Printf("Error navigating to URL: %v", err)
		os.Exit(1)
	}

	// アカウント入力
	elem, err := wd.FindElement(selenium.ByID, "mfid_user[email]")
	if err != nil {
		log.Printf("Error finding email field: %v", err)
		os.Exit(1)
	}
	err = elem.Clear()
	if err != nil {
		log.Printf("Error clearing email field: %v", err)
	}
	err = elem.SendKeys(user + string(selenium.EnterKey))
	if err != nil {
		log.Printf("Error entering email: %v", err)
		os.Exit(1)
	}

	// パスワード入力
	elem, err = wd.FindElement(selenium.ByID, "mfid_user[password]")
	if err != nil {
		log.Printf("Error finding password field: %v", err)
		os.Exit(1)
	}
	err = elem.Clear()
	if err != nil {
		log.Printf("Error clearing password field: %v", err)
	}
	err = elem.SendKeys(password + string(selenium.EnterKey))
	if err != nil {
		log.Printf("Error entering password: %v", err)
		os.Exit(1)
	}

	// CSVファイルを開く
	file, err := os.Open(inputFile)
	if err != nil {
		log.Printf("Error opening CSV file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	// CSVリーダーの作成
	reader := csv.NewReader(file)
	reader.LazyQuotes = true // クォートに関する厳密なチェックを無効化

	lineNum := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV line: %v", err)
			continue
		}

		lineNum++
		fmt.Printf("Start line[%d]\n", lineNum)

		// 日付が「#」の場合、コメント行として次へ飛ばす
		fmt.Printf("[0]%s\n", row[0])
		if row[0] == "#" || row[0] == "0" || row[0] == "計算対象" {
			fmt.Printf("[%d] Skip comment line!\n", lineNum)
			continue
		}

		// 「手入力」ボタンクリック
		elem, err = wd.FindElement(selenium.ByClassName, "cf-new-btn")
		if err != nil {
			log.Printf("Error finding 'manual input' button: %v", err)
			continue
		}
		err = elem.Click()
		if err != nil {
			log.Printf("Error clicking 'manual input' button: %v", err)
			continue
		}

		// 金額入力（収入/支出の切り替え）
		amount, err := strconv.Atoi(row[3])
		if err != nil {
			log.Printf("Error converting amount to integer: %v", err)
			continue
		}

		var plusMinusFlg string
		if amount > 0 {
			// 金額 > 0 ならば収入
			fmt.Printf("[%d] Plus!:\n", lineNum)
			fmt.Println(row)
			plusMinusFlg = "p"

			// 収入ボタンクリック
			elem, err = wd.FindElement(selenium.ByClassName, "plus-payment")
			if err != nil {
				log.Printf("Error finding 'plus' button: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking 'plus' button: %v", err)
				continue
			}
		} else if amount < 0 {
			// 金額 < 0 ならば支出
			fmt.Printf("[%d] Minus!:\n", lineNum)
			fmt.Println(row)
			amount *= -1 // 金額を正の値に変換
			plusMinusFlg = "m"
		} else {
			fmt.Printf("Error row num = %d\n", lineNum)
			continue
		}

		// 日付（YYYY/MM/DD）入力
		elem, err = wd.FindElement(selenium.ByID, "updated-at")
		if err != nil {
			log.Printf("Error finding date field: %v", err)
			continue
		}
		err = elem.Clear()
		if err != nil {
			log.Printf("Error clearing date field: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
		err = elem.SendKeys(row[1])
		if err != nil {
			log.Printf("Error entering date: %v", err)
			continue
		}
		err = elem.Click()
		if err != nil {
			log.Printf("Error clicking date field (1): %v", err)
		}
		err = elem.Click()
		if err != nil {
			log.Printf("Error clicking date field (2): %v", err)
		}
		time.Sleep(500 * time.Millisecond)

		// 金額入力
		elem, err = wd.FindElement(selenium.ByID, "appendedPrependedInput")
		if err != nil {
			log.Printf("Error finding amount field: %v", err)
			continue
		}
		err = elem.Clear()
		if err != nil {
			log.Printf("Error clearing amount field: %v", err)
		}
		err = elem.SendKeys(strconv.Itoa(amount))
		if err != nil {
			log.Printf("Error entering amount: %v", err)
			continue
		}

		// 大項目選択
		if row[5] != "未分類" {
			elem, err = wd.FindElement(selenium.ByID, "js-large-category-selected")
			if err != nil {
				log.Printf("Error finding large category field: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking large category field: %v", err)
				continue
			}

			elem, err = wd.FindElement(selenium.ByLinkText, row[5])
			if err != nil {
				log.Printf("Error finding large category option: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking large category option: %v", err)
				continue
			}
		}

		// 中項目選択
		if row[6] != "未分類" {
			subCategory := row[6]
			if len(subCategory) > 0 && subCategory[0] == '\'' {
				subCategory = subCategory[1:] // 先頭の「'」を削除
			}

			fmt.Printf("sub_category:%s\n", subCategory)
			elem, err = wd.FindElement(selenium.ByID, "js-middle-category-selected")
			if err != nil {
				log.Printf("Error finding middle category field: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking middle category field: %v", err)
				continue
			}

			elem, err = wd.FindElement(selenium.ByLinkText, subCategory)
			if err != nil {
				log.Printf("Error finding middle category option: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking middle category option: %v", err)
				continue
			}
		}

		// 内容入力
		var content string
		if row[7] == "" {
			content = row[2]
		} else {
			content = row[2] + "（" + row[7] + "）"
		}

		// 内容を50文字に制限
		if len(content) > 50 {
			content = content[:50]
		}

		elem, err = wd.FindElement(selenium.ByID, "js-content-field")
		if err != nil {
			log.Printf("Error finding content field: %v", err)
			continue
		}
		err = elem.Clear()
		if err != nil {
			log.Printf("Error clearing content field: %v", err)
		}
		err = elem.SendKeys(content)
		if err != nil {
			log.Printf("Error entering content: %v", err)
			continue
		}

		// 保存せずに閉じる（テストモード）
		time.Sleep(3 * time.Second)
		elem, err = wd.FindElement(selenium.ByClassName, "close")
		if err != nil {
			log.Printf("Error finding close button: %v", err)
			continue
		}
		err = elem.Click()
		if err != nil {
			log.Printf("Error clicking close button: %v", err)
			continue
		}
		time.Sleep(5 * time.Second)

		// 実際に保存する場合は以下を使用（現在はコメントアウト）
		/*
			time.Sleep(1 * time.Second)
			elem, err = wd.FindElement(selenium.ByID, "submit-button")
			if err != nil {
				log.Printf("Error finding submit button: %v", err)
				continue
			}
			err = elem.Click()
			if err != nil {
				log.Printf("Error clicking submit button: %v", err)
				continue
			}
		*/
	}

	fmt.Println("End :" + inputFile)
}
