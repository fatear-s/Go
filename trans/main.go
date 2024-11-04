package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type DictRequest struct {
	Query          string
	From           string
	To             string
	Reference      string
	CorpusIds      []string
	NeedPhonetic   bool
	Domain         string
	MilliTimestamp int64
}
type TooLTT struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   Data   `json:"data"`
}
type List struct {
	ID       string `json:"id"`
	ParaIdx  int    `json:"paraIdx"`
	Src      string `json:"src"`
	Dst      string `json:"dst"`
	Metadata string `json:"metadata"`
}
type Data struct {
	Event   string `json:"event"`
	Message string `json:"message"`
	List    []List `json:"list"`
}

func query(word string) {
	client := &http.Client{}
	request := DictRequest{Query: word, From: "en", To: "zh", Reference: "", CorpusIds: []string{}, NeedPhonetic: true, Domain: "common", MilliTimestamp: 1730524945975}
	buf, err := json.Marshal(request)
	if err != nil {
		return
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/ait/text/translate", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Acs-Token", "1730452072381_1730524945900_my6kzVABq17nl3mmJNqHPmlf1LCAcqPNTIqo49pwySAEsfqLtQbYOb7/gGkA0SGRRWQoqjbznmz2zTnF0qkFo+HhOu+ZSuUS9xNa2pqJJqv6KwoU0XOYafyk8UFXA6GtBJL+Q9geAervqVTqmhTNiFQwOIC3AvrUDNNSKfE/jM3Vz7aQJld2XXVU4H4mb2n3CwTuTFxvzu7vb+6G7pi/9fYqAFVytIloDd3hvLgAEBJ8YISoF/YMYoHqURhgzKXUDyDSKyO52rOghY7iUjlKFtMyV3M/0QLVfuDZK3x4euEBYOx/CxVbDbHREfznZ6s5qmsI8NuqxNKHjlBsuDJaLa42Ty+qMwB8mqwiIA+ZAr5+XHz45qOgU7AEKr5GbPnoVzeDZHK2pydH5JMrWRaTf/UGYf/xVkTK1vbvCFuygTS0foeBo5TxfQe51PBOtUXH5nDEztR5fyaQUdogEtZtsHQ6vZIfXd0bZBrY1mc3Xlp05fC4bAXF9DdaLCyDwrxs")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", `BIDUPSID=32023423FD2BD496DD1E34EABAA73E70; PSTM=1702628081; BAIDUID=32023423FD2BD496C5AFCE7747E7519D:FG=1; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1702907220; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; H_WISE_SIDS_BFESS=60272_60852_60887_60875_60897; H_WISE_SIDS=60272_60852_60897_60949_60980; BDSFRCVID=5skOJexroG3Wd8QJKWMWuFMPhmKK0gOTDYLEUamaI2AU2V4VY-HoEG0Pt_U-mEt-J8jwogKKLmOTHpKF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tbkD_C-MfIvDqTrP-trf5DCShUFsLboCB2Q-XPoO3KOReKnhyxRYyML0Q4cTaPriWbRM2MbgylRp8P3y0bb2DUA1y4vpXCr7a2TxoUJ2XMKVDq5mqfCWMR-ebPRiWPr9QgbjahQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hI0ljj82e5P0hxry2Dr2aI52B5r_5TrjDnCrjMOOXUI82h5y05JU5nbIaqvSMJohsl6MM-cvyT8sXnORXx74QC6L3lI-0xKKSCbKWDc1hUL1Db3JyhLLamTJslFy2t3oepvoX-cc3MkQyPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtbkD_C-MfIvhDRTvhCcjh-FSMgTBKI62aKDsMbDMBhcqJ-ovQTbRb5kuy4R3Ltrt3ejTBn6cWKJJ8UbeWfvp3t_D-tuH3lLHQJnph66dah5nhMJmBp_VhfL3qtCOaJby523i2n6vQpn2OpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0DTbLjH8jqTna--oa3RTeb6rjDnCrKT7dXUI82h5y05JU5nbIaqkbtqohsl6MM-cvyT8sXnORXx74B5vvbPOMthRnOlRKWDc1hUL1Db3JyhLLamTJslFy2t3oepvoX-cc3MkQyPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtbkD_C-MfIvhDRTvhCcjh-FSMgTBKI62aKDs_4OoBhcqJ-ovQTbRb5kuy4RUWlbt3ejTBn6cWKJJ8UbeWfvp3t_D-tuH3lLHQJnph66dah5nhMJmBp_VhfL3qtCOaJby523i2n6vQpn2OpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0-nDSHH_DJT-83j; BAIDUID_BFESS=32023423FD2BD496C5AFCE7747E7519D:FG=1; BDSFRCVID_BFESS=5skOJexroG3Wd8QJKWMWuFMPhmKK0gOTDYLEUamaI2AU2V4VY-HoEG0Pt_U-mEt-J8jwogKKLmOTHpKF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tbkD_C-MfIvDqTrP-trf5DCShUFsLboCB2Q-XPoO3KOReKnhyxRYyML0Q4cTaPriWbRM2MbgylRp8P3y0bb2DUA1y4vpXCr7a2TxoUJ2XMKVDq5mqfCWMR-ebPRiWPr9QgbjahQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hI0ljj82e5P0hxry2Dr2aI52B5r_5TrjDnCrjMOOXUI82h5y05JU5nbIaqvSMJohsl6MM-cvyT8sXnORXx74QC6L3lI-0xKKSCbKWDc1hUL1Db3JyhLLamTJslFy2t3oepvoX-cc3MkQyPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtbkD_C-MfIvhDRTvhCcjh-FSMgTBKI62aKDsMbDMBhcqJ-ovQTbRb5kuy4R3Ltrt3ejTBn6cWKJJ8UbeWfvp3t_D-tuH3lLHQJnph66dah5nhMJmBp_VhfL3qtCOaJby523i2n6vQpn2OpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0DTbLjH8jqTna--oa3RTeb6rjDnCrKT7dXUI82h5y05JU5nbIaqkbtqohsl6MM-cvyT8sXnORXx74B5vvbPOMthRnOlRKWDc1hUL1Db3JyhLLamTJslFy2t3oepvoX-cc3MkQyPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtbkD_C-MfIvhDRTvhCcjh-FSMgTBKI62aKDs_4OoBhcqJ-ovQTbRb5kuy4RUWlbt3ejTBn6cWKJJ8UbeWfvp3t_D-tuH3lLHQJnph66dah5nhMJmBp_VhfL3qtCOaJby523i2n6vQpn2OpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0-nDSHH_DJT-83j; BDRCVFR[Zh1eoDf3ZW3]=mk3SLVN4HKm; delPer=0; ZFY=3qMALeR0UI0YbQfifCImgkVucnqXfLvmypesbnhPNRQ:C; PSINO=6; BA_HECTOR=80ak810h8h0l01a50000858g8ntvd51jib2ne1u; BDORZ=FFFB88E999055A3F8A630C64834BD6D0; H_PS_PSSID=60272_60852_60949_60980_60941_61027_61036; log_first_time=1730524016115; log_last_time=1730524016115; ab_sr=1.0.1_NmZmMTFmMmI5YTc4N2Q5ZGNjNTNjYmFmMzNkYTdjNjI4YjI1ZjYyZjI2MWY0MGQ2MmIxZTNmNmMyNGFlNTQ2YmUyOTY1YjI1NDMyZDllZjFjNTc4NzRkYzgwM2M1YTk3MWE2OTRmYjRkYjhhZWIxMzQ3MTIwYWE3OWZlYzRlMTlkYjVjMDBmNDk3MTE4NWYwZjM2ZDBjYzFkZGFlODI4OGY2OTM4MzZlMjljMjM2NmJhM2ViNjQ1YjcxZmYyMmFj; RT="z=1&dm=baidu.com&si=3648d1bd-520a-4701-85c7-10ad857ae80f&ss=m2zpbhbg&sl=3&tt=3ni&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=jqyc"`)
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/mtpe-individual/multimodal?query=good&lang=en2zh&ext_channel=Aldtype")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/event-stream")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := bytes.Split(bodyText, []byte("event: message\ndata:"))
	var result_Dict TooLTT
	for _, msg := range result {
		err := json.Unmarshal(msg, &result_Dict)
		if err != nil {
			//fmt.Println("error")
			continue
		}
	}
	fmt.Println(word, ":", result_Dict.Data.List[0].Dst)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	query(word)
}
