package function

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

var verificationToken string

// グローバル変数を定義して、init で取得処理を書くと次回の関数呼び出し時にも再利用される。
// https://cloud.google.com/functions/docs/bestpractices/tips#use_global_variables_to_reuse_objects_in_future_invocations
func init() {
	// Cloud Function に環境変数として登録済み
	verificationToken = os.Getenv("VERIFICATION_TOKEN")
}

// エントリポイントとして渡す関数
func HelloCommand(w http.ResponseWriter, r *http.Request) {
	// スラッシュコマンドのリクエストをパースする。
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Slack から来るリクエストに付与される ValidateToken をチェックする。
	if !s.ValidateToken(verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/hello":
		// ResponseType のデフォルトは ephemeral になっており、ephemeral　では、投稿者にしかメッセージが表示されない。
		// チャンネル全体に投稿する時は、in_channel を指定する。
		params := &slack.Msg{
			ResponseType: "in_channel",
			Text:         "こんにちは、<@" + s.UserID + ">さん",
		}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
