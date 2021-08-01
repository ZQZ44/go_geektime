import (
        "database/sql"
        "os"
        "fmt"
        "github.com/pkg/errors"
)

/*  在 dao 层遇到sql.ErrNoRows，应该调用Warp 及时定位错误并保存堆栈信息，然而返回给上层，防止无法定位错误发生点的麻烦。
*/

func main() {
    err := dao()
    if err != nil {
        fmt.Printf("Error: %+v\n", err)
        os.Exit(1)
    }
}

func dao() err {
    db, err := sql.open("......",".....")
    var sth string
    err = db.QueryRow("......").Scan(&sth)
    if err == sql.ErrNoRows {
        return errors.Warp(err, "dao sql ErrNoRow")
    }
}
