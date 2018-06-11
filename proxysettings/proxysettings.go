package proxysettings

import "C"
//#cgo LDFLAGS: -lwininet
//#include <stdlib.h>
//#include <windef.h>
//#include <winbase.h>
//#include <wininet.h>
//int clear() {
//  INTERNET_PER_CONN_OPTION_LIST options;
//  options.dwOptionCount = 3;
//  options.pOptions = (INTERNET_PER_CONN_OPTION*)calloc(options.dwOptionCount, sizeof(INTERNET_PER_CONN_OPTION));
//  options.pOptions[0].dwOption = INTERNET_PER_CONN_FLAGS;
//  options.pOptions[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER;
//  options.pOptions[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS;
//  options.pOptions[2].Value.pszValue = TEXT("<local>");
//  options.pszConnection = NULL;
//	options.pOptions[0].Value.dwValue = PROXY_TYPE_DIRECT;
//  InternetSetOption(NULL, INTERNET_OPTION_PER_CONNECTION_OPTION, &options, sizeof(INTERNET_PER_CONN_OPTION_LIST));
//	free(options.pOptions);
//}
//int set(char* addr) {
//  INTERNET_PER_CONN_OPTION_LIST options;
//  options.dwOptionCount = 3;
//  options.pOptions = (INTERNET_PER_CONN_OPTION*)calloc(options.dwOptionCount, sizeof(INTERNET_PER_CONN_OPTION));
//  options.pOptions[0].dwOption = INTERNET_PER_CONN_FLAGS;
//  options.pOptions[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER;
//  options.pOptions[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS;
//  options.pOptions[2].Value.pszValue = TEXT("<local>");
//  options.pszConnection = NULL;
//	options.pOptions[0].Value.dwValue = PROXY_TYPE_AUTO_PROXY_URL;
//	options.pOptions[1].dwOption = INTERNET_PER_CONN_AUTOCONFIG_URL;
//	options.pOptions[1].Value.pszValue = addr;
//  InternetSetOption(NULL, INTERNET_OPTION_PER_CONNECTION_OPTION, &options, sizeof(INTERNET_PER_CONN_OPTION_LIST));
//	free(options.pOptions);
//}
import "C"

func Clear() {
	C.clear()
}

func Set(pacAddr string) {
	C.set(C.CString(pacAddr))
}
