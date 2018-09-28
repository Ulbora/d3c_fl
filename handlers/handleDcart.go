/*
 Copyright (C) 2018 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.

 Copyright (C) 2018 Ken Williamson
 All rights reserved.

 Certain inventions and disclosures in this file may be claimed within
 patents owned or patent applications filed by Ulbora Labs LLC., or third
 parties.

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published
 by the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package handlers

import (
	dcd "FindByFlic/dbdelegate"
	"encoding/json"
	"io/ioutil"
	//dbi "github.com/Ulbora/dbinterface"
	"log"
	"net/http"
)

//HandleDcartIndex HandleDcartIndex
func (h *Handler) HandleDcartIndex(w http.ResponseWriter, r *http.Request) {
	h.Sess.InitSessionStore(w, r)
	var order = r.URL.Query().Get("order")
	var carturl = r.URL.Query().Get("carturl")
	session := h.getSession(w, r)
	session.Values["order"] = order
	session.Values["carturl"] = carturl
	serr := session.Save(r, w)
	if serr != nil {
		log.Println("Session Err:", serr)
	}

	// 	SecureUrl: 3dcart merchant's Secure URL.
	// PrivateKey: Your application's private key.
	// Token: The 3dcart merchant's token.
	//<iframe src="https://localhost:8070?cart=3dcar&carturl=[store_url]"></iframe>
	//<iframe src="https://localhost:8070?cart=3dcar&carturl=[store_url]"></iframe>

	// secureURL := r.Header.Get("SecureUrl")
	// privateKey := r.Header.Get("PrivateKey")
	// token := r.Header.Get("Token")
	var p PageParams
	p.Order = order
	p.CartURL = carturl
	// p.URL = secureURL
	// p.Key = privateKey
	// p.Token = token
	h.Templates.ExecuteTemplate(w, "dcartIndex.html", &p)
}

//HandleDcartConfig HandleDcartConfig
func (h *Handler) HandleDcartConfig(w http.ResponseWriter, r *http.Request) {
	h.Sess.InitSessionStore(w, r)
	// var cart = r.URL.Query().Get("cart")
	// var carturl = r.URL.Query().Get("carturl")
	// session := h.getSession(w, r)
	// session.Values["cart"] = cart
	// session.Values["carturl"] = carturl
	// serr := session.Save(r, w)
	// log.Println(serr)
	// 	SecureUrl: 3dcart merchant's Secure URL.
	// PrivateKey: Your application's private key.
	// Token: The 3dcart merchant's token.
	//<iframe src="https://localhost:8070?cart=3dcar&carturl=[store_url]"></iframe>
	//<iframe src="https://localhost:8070?cart=3dcar&carturl=[store_url]"></iframe>

	// secureURL := r.Header.Get("SecureUrl")
	// privateKey := r.Header.Get("PrivateKey")
	// token := r.Header.Get("Token")
	// var p PageParams
	// p.Cart = cart
	// p.CartURL = carturl
	// p.URL = secureURL
	// p.Key = privateKey
	// p.Token = token
	h.Templates.ExecuteTemplate(w, "dcartConfig.html", nil)
}

//HandleDcartCb HandleDcartCb
func (h *Handler) HandleDcartCb(w http.ResponseWriter, r *http.Request) {
	var rtn bool
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("body err: ", err)
	} else {
		log.Println("json from dcart at start:", string(b))
	}

	log.Println("service call from 3dcart on callback-------")
	cType := r.Header.Get("Content-Type")
	if cType != "application/json" {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		//var dcReg dcd.DCartUser
		//log.Println("body from dcart :", r.Body)
		dcReg := new(dcd.DCartUser)
		json.Unmarshal(b, dcReg)
		// decoder := json.NewDecoder(r.Body)
		// err := decoder.Decode(&dcReg)
		// if err != nil {
		// 	log.Println("Decode error: ", err.Error())
		// 	//http.Error(w, err.Error(), http.StatusBadRequest)
		// }

		if dcReg.Action == "AUTHORIZE" {
			rtn, _ = h.FindFFLDCart.AddUser(dcReg)
		} else if dcReg.Action == "REMOVE" {
			rtn = h.FindFFLDCart.RemoveUser(dcReg)
		}

		jsn, err := json.Marshal(dcReg)
		if err != nil {
			log.Println("marshal err: ", err)
		}
		log.Println("json from dcart at end:", string(jsn))
		if rtn {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

//HandleDcartFindFFL HandleDcartFindFFL
func (h *Handler) HandleDcartFindFFL(w http.ResponseWriter, r *http.Request) {
	zip := r.FormValue("zip")
	res := h.FFLFinder.FindFFL(zip)
	var pg FFLPageParams
	pg.FFLList = res
	h.Templates.ExecuteTemplate(w, "dcartAddFfl.html", &pg)
}
