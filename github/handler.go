package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/wegotour/webhooks"
	"github.com/wegotour/webhooks/gcf"
	"github.com/whatsauth/wa"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)
	link := "https://medium.com/@daffaaudyapramana03/cara-menggunakan-whatsauth-free-2fa-otp-notif-whatsapp-gateway-api-gratis-458105965936"
	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		if msg.Message == "loc" || msg.Message == "Loc" || msg.Message == "lokasi" || msg.LiveLoc {
			location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
			if err != nil {
				// Handle the error (e.g., log it) and set a default location name
				location = "Unknown Location"
			}

			reply := fmt.Sprintf("Aku ramal kamu pasti berada di %s \n Koordinatnya : %s - %s\n Tutorial WhatsAuth Ada di link dibawah ini"+
				" yaa %s\n", location,
				strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)), link)
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: reply,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		
			} else if msg.Message == "Mohon" || msg.Message == "Untukk" || msg.Message == "Bersabar" {
				dt := &wa.TextMessage{
					To:       msg.Phone_number,
					IsGroup:  false,
					Messages: fmt.Sprintf("Tolong jangan spam %s akan segera kami proses", msg.Alias_name),
				}
				resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
	
			} else if msg.Message == "Pak" || msg.Message == "Bu" || msg.Message == "Kak" {
				dt := &wa.TextMessage{
					To:       msg.Phone_number,
					IsGroup:  false,
					Messages: fmt.Sprintf("Terima kasih %s silahkan datang kembali", msg.Alias_name),
				}
				resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
	
			} else if strings.Contains(msg.Message, "login") {
				//login username test password testcihuy
				messages := strings.Split(msg.Message, " ")
				email := messages[2]
				password := messages[len(messages)-1]
				dt := &webhooks.Logindata{
					Email:    email,
					Password: password,
				}
				res, _ := atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://asia-southeast2-pasabar.cloudfunctions.net/Admin-Login")
				dat := &wa.TextMessage{
					To:       msg.Phone_number,
					IsGroup:  false,
					Messages: res.Response,
				}
				resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dat, "https://api.wa.my.id/api/send/message/text")
			} else {
				randm := []string{
					"Hi" + msg.Alias_name + "\n Daffa & Prisya lagi sibuk \n aku bot wegotour salam kenal yaa \n Cara penggunaan WhatsAuth ada di link berikut ini ya kak...\n" + link,
					"Mohon jangan spam kak",
					"Untuk pemesanan kita sedang ada gangguan",
					"Kita dalam proses perbaikan, mohon ditunggu",
					"Kita akan segera kembali, mohon bersabar ya kak",
				}
				dt := &wa.TextMessage{
					To:       msg.Phone_number,
					IsGroup:  false,
					Messages: function.GetRandomString(randm),
				}
				resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
			}
		} else {
			resp.Response = "Secret Salah"
		}
		fmt.Fprintf(w, resp.Response)
	}