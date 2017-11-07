#pragma once
#include "windows.h"
#include "stdio.h"

HWND foregroundWindow();
WCHAR * windowText(HWND);
int wcharlen(WCHAR * str);
int listWindows(HWND *);

HWND foregroundWindow() {
    HWND window = GetForegroundWindow();
    return window;
}

WCHAR * windowText(HWND window) {
    WCHAR * title = (WCHAR *)malloc(sizeof(WCHAR) * 1024+1);
    int len = GetWindowTextW(window, title, 1024);
    title[len+1] = '\0';

    return title;
}

int wcharlen(WCHAR * str) {
    return wcslen(str);
}

