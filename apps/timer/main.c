#include "raylib.h"
#include <time.h>

int main() {
  InitWindow(2560, 1440, "timer");
  SetConfigFlags(FLAG_WINDOW_RESIZABLE | FLAG_WINDOW_MAXIMIZED);

  while (!WindowShouldClose()) {
    time_t cur_time;
    time(&cur_time);

    struct tm target_tm = { 0, 39, 11, 11, 12 - 1, 2025 - 1900, .tm_isdst = -1, };
    time_t target_time = mktime(&target_tm);

    double diff = difftime(target_time, cur_time);

    int mins = (int)diff / 60;
    int secs = (int)diff % 60;
  }
}
