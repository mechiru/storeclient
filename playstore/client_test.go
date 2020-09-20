package playstore

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetURL(t *testing.T) {
	for idx, c := range []struct {
		in   string
		want string
		err  error
	}{
		{"a.b.c", "https://play.google.com/store/apps/details?hl=ja&id=a.b.c", nil},
		{"", "", errEmptyBundleID},
	} {
		url, err := NewClient(Lang("ja")).getURL(c.in)
		if err != c.err {
			t.Errorf("idx=%d: got=%+v, want=%+v", idx, err, c.err)
		}
		if err != nil {
			continue
		}
		got := url.String()
		if got != c.want {
			t.Errorf("idx=%d: got=%s, want=%s", idx, got, c.want)
		}
	}
}

func TestGet(t *testing.T) {
	if _, err := NewClient().Get(context.Background(), "com.cookpad.android.activities"); err != nil {
		t.Error(err)
	}
}

// https://play.google.com/store/apps/details?hl=ja&id=com.cookpad.android.activities
func TestParseHTML(t *testing.T) {
	c, err := ioutil.ReadFile("./testdata/com.cookpad.android.activities.html")
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseHTML(bytes.NewReader(c))
	if err != nil {
		t.Fatal(err)
	}

	want := &Detail{
		Title: "クックパッド-No.1料理レシピ検索アプリ",
		Description: `月次利用者数約5800万人・掲載レシピ数320万品以上！
日本最大の料理レシピ検索・投稿サービス「クックパッド」の公式アプリ。
レシピに特化して、最大限使いやすく・便利にしたアプリでは、時短、節約、簡単、ダイエットなど気になるテーマからお菓子、パンなどのカテゴリー、料理名や食材から検索できて、どんなレシピも簡単に見つかります。

「つくれぽ」で作ったみんなの感想も分かる！
気になったレシピはMYフォルダに保存、作った料理を記録してあなただけのレシピ帳に！

さらにプレミアムサービス会員なら、人気順検索や、殿堂入りレシピの機能でみんなに人気のレシピを快適に素早く探せます！


■こんな方におすすめ
・みんなが作っている料理を知りたい
・冷蔵庫に今ある食材で料理したい
・時短でおいしく料理したい
・節約しながら、家族が喜ぶ料理をつくりたい
・マンネリを脱出し、レパートリーを増やしたい
・毎日の献立に使える簡単レシピが知りたい
・旬の食材をおいしく食べたい
・毎日のお弁当おかずのマンネリを脱したい


■クックパッドの特徴
＜どんなレシピも見つかる＞
320万品以上あるレシピから、お肉、お魚、お野菜、お弁当やお菓子、パンなどのカテゴリ、和食や中華などのジャンル、ハンバーグやカレーといった料理名や食材から検索できて、どんなレシピも簡単に見つかります。

＜家庭で作りやすいレシピ＞
日本中のみんなが投稿したレシピが320万品以上。普段の献立に使える簡単レシピから、本格派のこだわりレシピまで、作りやすくておいしいレシピがそろっています。

＜今日作られている料理がわかる＞
みんなに話題のおすすめレシピは毎日更新。新機能「タイムライン」では、今みんなに作られているレシピがリアルタイムにわかります。

＜あなただけのレシピ帳に＞
作りたいレシピは20品まで保存可能。スマホで撮った料理写真を自動で整理してくれる「料理きろく」も加わって、ますます便利になりました。


■さらに！プレミアムサービス（月額302円（税込））でできること
・人気順検索：320万品以上の中から人気のレシピがすぐわかる
・殿堂入りレシピ：失敗知らずの鉄板レシピ集
・デイリーアクセス数ランキング：アクセスが多かったレシピを日替わりで紹介
・プレミアム献立：管理栄養士監修の献立を毎日提案
・専門家厳選レシピ：ダイエットや離乳食などプロ監修の目的別レシピ
・レシピ保存：MYフォルダの容量が3,000件に大幅UP`,
		CoverArtURL:   "https://lh3.googleusercontent.com/G8KxoLSAJIYrDSObg07KNi55XAij9uO4hr4VQYXTTmfCpvfHswL0SwnGgx3Zcfvj2Hk=s180",
		ContentRating: "3 歳以上",
		GenreID:       "FOOD_AND_DRINK",
		Genre:         "フード＆ドリンク",
		DeveloperID:   "6698899737769238815",
		Developer:     "Cookpad Inc.",
		DeveloperURL:  "https://cookpad.com/",
		BundleID:      "com.cookpad.android.activities",
		StoreID:       "com.cookpad.android.activities",
	}

	if !cmp.Equal(got, want) {
		t.Errorf(cmp.Diff(got, want))
	}
}
