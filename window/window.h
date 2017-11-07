#pragma once
#include "windows.h"
#include "stdio.h"

typedef struct ListOfHandles ListOfHandles_t;
HWND foregroundWindow();
WCHAR * windowText(HWND);
int wcharlen(WCHAR * str);
ListOfHandles_t getAllWindows();
BOOL CALLBACK getAllWindowsProc(HWND handle, LPARAM lParam);

HWND foregroundWindow() {
    HWND window = GetForegroundWindow();
    return window;
}

typedef struct ListOfHandles {
    int count;
    int size;
    HWND * handles;
} ListOfHandles_t;


ListOfHandles_t getAllWindows() {
    ListOfHandles_t handles;
    handles.count   = 0;
    handles.size    = 64;
    handles.handles = (HWND *)malloc(sizeof(HWND) * 64);
    EnumWindows(getAllWindowsProc, (LPARAM)&handles);
    return handles;
}

BOOL CALLBACK getAllWindowsProc(HWND handle, LPARAM lParam) {
    ListOfHandles_t * handles = (ListOfHandles_t *)lParam;
    // Increase the size of array when needed.
    if ( handles->count >= handles->size ) {
        handles->size *= 2;
        handles->handles = (HWND *)realloc(handles->handles, sizeof(HWND) * handles->size);
    }
    handles->handles[handles->count] = handle;
    handles->count++;

    return 1;
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

