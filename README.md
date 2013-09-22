Go語言在GAE上操作CRUD基本範例 
--------------------------------------------------
GAE(Google App Engine)平台上使用Go語言操作CRUD的範例程式碼, 資料庫使用GAE的Datastore。

##前言

Go語言是Google開發的編譯型語言, 於2009年推出, 關於Go的簡介在此不贅述, 請參考[Go Wike](http://zh.wikipedia.org/wiki/Go)。

Go語言在網路上的學習資源還不是很多, 尤其是在GAE的應用更少, 此範例幫助對Go語言有興趣的新手了解並學習Go語言的程式形式, 同時對GAE及其Datastore有基本的認識。

GAE是「Google應用服務引擎」在Google的平台上執行您的網路應用程式, 有**免費**的配額可以使用, 起初僅支援Python、JAVA, 後來才加入了Go、PHP等。

Datastore是GAE是的Database與一般RDBMS有些差異, 本人我也還不是很瞭解, 網路上的參考資料也不是很多, 簡單想就是Key-Value, 這些Key甚至可以有Parent, 實作過程發現常會使用到Key, datastore的方法裡若要對資料操作也常需要Key作為引數。

-  _注意: Go語言的主要應用並不是在網頁上, 只是此範例從網頁應用切入學習Go語言_。

## Demo

* [http://golang-crud.appspot.com/](http://golang-crud.appspot.com/)

## Getting Started


1. 申請 [Google App Engine](http://appengine.google.com), Create Application, 請記下您的**Application Identifier**。

2. 安裝Go Tools: [http://golang.org/doc/install](http://golang.org/doc/install)。

3. 安裝[App Engine SDK](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)。

4. Clone:
	
		git clone git@github.com:yfxie/go-lang-crud-on-gae.git
		cd go-lang-crud-on-gae/

	建立app.yaml檔案, 內容如下: _application請填入Step1的Application Identifier_

		application: helloworld
		version: 1
		runtime: go
		api_version: go1

		handlers:
		- url: /.*
		  script: _go_app

5. Local Testing: 
		
		dev_appserver.py .

	Server: http://localhost:8080

	AdminServer: http://localhost:8000

6. Deploy to GAE:

		appcfg.py update .

	Your app will running at _http://yourname.appspot.com_

## Package

- **appengine**
	
	用來產生context的必要包, 與GAE相關的函數經常須使用context做為引數。

- **appengine/datastore**
	
	操作GAE Datastore的工具。

- **appengine/user**

	輕鬆實現Google使用者驗證。

- **fmt**
	
	純文字輸出用。

- **html/template**
	
	內建樣板引擎。

- **net/http**

	Server建路由。

- **strconv**
	
	字串形態轉換用。

- **time**
	
## Tools

* [GoSublime](https://github.com/DisposaBoy/GoSublime): 好用的Sublime的Go套件, 自動完成、程式碼整理功能..etc。

## Reference

- [Author's blog article](http://blog.developer.tw/post/621-learning-go-language-operating-crud-on-gae)
- [GoLang.org](http://golang.org/)
- [DataStoreReference - GAE](https://developers.google.com/appengine/docs/go/datastore/reference)
