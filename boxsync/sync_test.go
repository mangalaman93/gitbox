package boxsync

import (
	"testing"
	"fmt"
)

//func TestDownloadFile(t *testing.T) {
//	Download_File(94149725862, "/tmp/jane.jpeg")
//}

//func TestGet_Folder_Id(t *testing.T) {
//	Get_Folder_Id("home/pratyaksh/abc/")
//}

//func TestGet_Folder_Items(t *testing.T) {
//	Get_Folder_Items(0)
//}
//

//func TestDownload_File_By_Path(t *testing.T) {
//	Download_File_By_Path("L1/L2/L3/file1.txt", "/tmp/file1.txt")
//}

//func TestGet_Object_Id(t *testing.T) {
//	s, ty, err := Get_Object_Id("L1/L2/L3/file1.txt")
//	fmt.Println(s)
//	fmt.Println(ty)
//	fmt.Println(err)
//}

func TestDownload_Folder(t *testing.T) {
	err := Download_Folder("f1", "/tmp/")
	fmt.Println(err)
}
