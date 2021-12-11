//@Title		db.go
//@Description	数据库相关操作
//@Author		zy
//@Update		2021.12.5

package common

import (
	"csa_6/userinformation"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// DB 定义一个全局变量
var DB *sql.DB

//InitDB
//@title		InitDB()
//@description	连接数据库
//@author		zy
//@param		dsn string
//@return		*sql.DB error
func InitDB(dsn string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("failed to open database, err:%v", err)
		return nil, err
	}
	err = DB.Ping()

	if err != nil {
		fmt.Printf("failed to connect database, err: %v", err)
		return nil, err
	}
	return DB, err
}

//QueryUserInfo
//@title		QueryUserInfo()
//@description	查询用户名u.Username是否已经存在
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QueryUserInfo(u userinformation.UserInfo) bool {
	sqlStr := "select username from user where username=?"			//sql语句
	var UTemp userinformation.UserInfo
	err := DB.QueryRow(sqlStr, u.Username).Scan(&UTemp.Username)	//调用QueryRow进行插入
	if err != nil {
		fmt.Printf("用户名%s还未注册\n", u.Username)
		return false
	}
	fmt.Printf("用户:%s你好\n", UTemp.Username)
	return true
}

//QueryUserInfoExist
//@title		QueryUserInfoExist()
//@description	查询用户名和密码是否正确
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QueryUserInfoExist(u userinformation.UserInfo) bool {
	sqlStr := "select username, password from user where username=? and password=?"	//sql语句
	var UTemp userinformation.UserInfo

	err := DB.QueryRow(sqlStr, u.Username, u.Password).Scan(&UTemp.Username, &UTemp.Password)	//调用QueryRow进行插入

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	fmt.Printf("用户:%s你好\n", UTemp.Username)
	return true
}

//InsertUserInfo
//@title		InsertUserInfo()
//@description	注册成功时插入到user表中
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func InsertUserInfo(u userinformation.UserInfo) bool{
	sqlStr := "insert into user(username, password, secret_key, secret_value) values (?,?,?,?)"	//sql语句
	ret, err := DB.Exec(sqlStr, u.Username, u.Password, u.SecretKey, u.SecretValue)				//插入操作
	if err != nil {
		fmt.Printf("insert failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}


//QuerySecretKeyValue
//@title		QuerySecretKeyValue()
//@description	查询用户米和密保是否正确
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QuerySecretKeyValue(u userinformation.UserInfo) bool{
	sqlStr := "select username, secret_key, secret_value from user where username=? and secret_key=? and secret_value=?"
	var uTemp userinformation.UserInfo

	err := DB.QueryRow(sqlStr, u.Username, u.SecretKey, u.SecretValue).Scan(&uTemp.Username, &uTemp.SecretKey, &uTemp.SecretValue)
	if err != nil {
		fmt.Printf("SecretKeyValue scan failed, err:%v\n", err)
		return false
	}
	fmt.Printf("用户:%s你好\n", uTemp.Username)
	return true
}


//AlterUserInfo
//@title		AlterUserInfo()
//@description	更改用户密码
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func AlterUserInfo(u userinformation.UserInfo) bool{
	sqlStr := "update user set password=? where username=?"
	_, err := DB.Exec(sqlStr, u.Password, u.Username)										//更新对应的值——favor++
	if err != nil {																		//存在错误 提示用户
		fmt.Printf("failed to update, err:%v\n", err)
		return false
	}
	fmt.Printf("success to update password\n")
	return true
}

//InsertArticle
//@title		InsertArticle()
//@description	成功发表文章时插入到blog表中
//@author		zy
//@param		uArticle userinformation.Article
//@return		bool
func InsertArticle(uArticle userinformation.Article) bool{
	sqlStr := "insert into blog(username, title, content) values (?,?,?)"			//sql语句
	ret, err :=DB.Exec(sqlStr, uArticle.Username, uArticle.Title, uArticle.Content)	//sql操作
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}

//DeleteArticle
//@title		DeleteArticle()
//@description	成功删除文章时将blog表中对用的数据删除
//@author		zy
//@param		uArticle userinformation.Article
//@return		bool
func DeleteArticle(uArticle userinformation.Article) bool{
	sqlStr := "delete from blog where username=? and title=?"			//sql语句
	ret, err := DB.Exec(sqlStr, uArticle.Username, uArticle.Title)		//sql操作
	if err != nil {
		fmt.Printf("delete fail, err: %v\n", err)
		return false
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("affect fail, err: %v\n", err)
		return false
	}
	if n == 0 {
		fmt.Printf("affect 0 line, err: %v\n", err)
		return false
	}
	fmt.Printf("delete success, delete %d article which title is %s", n, uArticle.Title)
	return true
}

//LikeArticle
//@title		LikeArticle()
//@description	点赞别人的文章
//@author		zy
//@param		uArticle userinformation.Article
//@return		bool
func LikeArticle(uArticle userinformation.Article) int{
	sqlStr := "update blog set favor=favor+1 where username=? and title = ?"			//sql语句
	ret, err := DB.Exec(sqlStr, uArticle.Username, uArticle.Title)										//更新对应的值——favor++
	if err != nil {																		//存在错误 提示用户
		fmt.Printf("failed to update, err:%v", err)
		return 1
	}
	n, err1 := ret.RowsAffected()														//判断更新了几行相应的值 即判断username和title是否存在
	if err1 != nil {
		fmt.Printf("failed to affect, err:%v", err)
		return 1
	}
	if n == 0 {																			//没有相应的博客
		fmt.Printf("faied to like because %s doesn't write %s , err:%v", uArticle.Username, uArticle.Title, err)
		return 2
	}

	fmt.Printf("success like")
	return 3
}

//MessageInsert
//@title		MessageInsert()
//@description	留言成功时插入msg表
//@author		zy
//@param		Msg userinformation.Msg
//@return		bool
func MessageInsert(Msg *userinformation.Msg) bool{
	// 判断是否有Msg.OtherUsername 这个用户
	if !QueryUserInfo(userinformation.UserInfo{Username: Msg.OtherUsername}) {
		fmt.Printf("没有这个用户")
		return false
	}
	sqlMax := "select MAX(msg_id) from msg"
	var i int
	err := DB.QueryRow(sqlMax).Scan(&i)
	if err != nil {
		fmt.Printf("Max query is wrong, err:%v\n", err)
	}
	Msg.MsgId = i+1
	sqlStr := "insert into msg(username, otherusername, message, msg_id, msg_time) values (?,?,?,?,?)"			//sql语句
	ret, err := DB.Exec(sqlStr, Msg.Username, Msg.OtherUsername, Msg.Message, Msg.MsgId, time.Now())				//插入操作
	if err != nil {
		fmt.Printf("insert failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}


//ReplyMessage
//@title		ReplyMessage()
//@description	输入otherusername message msgid回复留言
//@author		zy
//@param		Msg userinformation.Msg
//@return		bool
func ReplyMessage(Msg userinformation.Msg) bool{
	// 判断是否有Msg.OtherUsername 这个用户
	if !QueryUserInfo(userinformation.UserInfo{Username: Msg.OtherUsername}) {
		fmt.Printf("没有这个%s用户", Msg.OtherUsername)
		return false
	}

	sqlStr1 := "select msg_id from msg where msg_id=?"
	var i int
	err := DB.QueryRow(sqlStr1, Msg.MsgId).Scan(&i)
	if err != nil {
		fmt.Printf("No msg_id is %d", Msg.MsgId)
		return false
	}
	sqlStr2 := "insert into msg(username, otherusername, message, msg_id, is_reply_message, msg_time) values(?,?,?,?,?,?)"
	ret, err := DB.Exec(sqlStr2, Msg.Username, Msg.OtherUsername, Msg.Message, Msg.MsgId, 1, time.Now())
	if err != nil {
		fmt.Printf("insert msg_reply failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert msg_reply success, the id is %d.\n", id)
	return true
}


//MessageDelete
//@title		MessageDelete()
//@description	删除留言
//@author		zy
//@param		Msg userinformation.Msg
//@return		bool
func MessageDelete(Msg userinformation.Msg) bool {
	sqlStr := "delete from msg where username=? and otherusername=? and message=?"			//sql语句
	ret, err := DB.Exec(sqlStr, Msg.Username, Msg.OtherUsername, Msg.Message)		//sql操作
	if err != nil {
		fmt.Printf("delete fail, err: %v\n", err)
		return false
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("affect fail, err: %v\n", err)
		return false
	}
	if n == 0 {
		fmt.Printf("affect 0 line, err: %v\n", err)
		return false
	}
	fmt.Printf("delete success, delete %d message which username is %s and otherusername is %s and message is %s\n", n, Msg.Username, Msg.OtherUsername, Msg.Message)
	return true
}