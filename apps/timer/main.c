#include "raylib.h"
#include <stdio.h>
#include <time.h>
#include <sys/time.h>

#define FONT_SIZE 400

int main() {
  InitWindow(2560, 1440, "timer");
  SetConfigFlags(FLAG_WINDOW_RESIZABLE | FLAG_WINDOW_MAXIMIZED);

  char minbuf[100];
  char secbuf[100];
  char milibuf[100];

  while (!WindowShouldClose()) {
    struct timeval tv;
    gettimeofday(&tv, NULL);

    struct tm target_tm = { 0, 39, 11, 11, 12 - 1, 2025 - 1900, .tm_isdst = -1, };
    time_t target_time_sec = mktime(&target_tm);

    long long target_ms = (long long)target_time_sec * 1000;
    long long current_ms = (long long)tv.tv_sec * 1000 + tv.tv_usec / 1000;

    long long diff_ms = target_ms - current_ms;

    int mins = diff_ms / (1000 * 60);
    int secs = (diff_ms % (1000 * 60)) / 1000;
    int millis = diff_ms % 1000;
    
    // Ensure milliseconds are always positive and within [0, 999] range
    if (diff_ms < 0) {
        if (millis != 0) { // Only adjust if not already 0 for negative time
            millis = 1000 - millis;
            secs--;
        }
        if (secs < 0) {
            secs = 60 + secs;
            mins--;
        }
    }

    sprintf(minbuf, "%d", mins);
    sprintf(secbuf, ":%d", secs);
    sprintf(milibuf, ".%d", millis);

    int fullwidth = MeasureText("00:00.0000", FONT_SIZE);
    int minwidth = MeasureText("00", FONT_SIZE);
    int secwidth = MeasureText("00:00", FONT_SIZE);

    BeginDrawing();
    ClearBackground(BLACK);
    DrawText(minbuf, GetScreenWidth() / 2 - fullwidth / 2, GetScreenHeight() / 2 - FONT_SIZE / 2, FONT_SIZE, WHITE);
    DrawText(secbuf, GetScreenWidth() / 2 - fullwidth / 2 + minwidth, GetScreenHeight() / 2 - FONT_SIZE / 2, FONT_SIZE, WHITE);
    DrawText(milibuf, GetScreenWidth() / 2 - fullwidth / 2 + secwidth, GetScreenHeight() / 2 - FONT_SIZE / 2, FONT_SIZE, WHITE);
    EndDrawing();
  }
}
