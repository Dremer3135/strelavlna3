import type { AdminViewedTeam } from "$lib/types";
import { writable } from "svelte/store";

export let adminViewedTeams = writable<Record<string, AdminViewedTeam>>({});