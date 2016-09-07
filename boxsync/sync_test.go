package boxsync

import (
	"testing"
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
//	s, ty, err := Get_Object_Id("f1/file2.txt")
//	fmt.Println(s)
//	fmt.Println(ty)
//	fmt.Println(err)
//}

//func TestDownload_Folder(t *testing.T) {
//	err := Download_Folder("f1", "/tmp/")
//	fmt.Println(err)
//}

//func TestUpload_File_By_Id(t *testing.T) {
//	Upload_File_By_Id("94153701752", "/tmp/f1/file2.txt")
//}

func TestUpload_New_File(t *testing.T) {
	Upload_New_File("0", "/tmp/f3.txt")
}
