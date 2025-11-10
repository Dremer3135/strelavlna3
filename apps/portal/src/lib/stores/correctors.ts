// src/lib/stores/correctors.ts
import { writable } from 'svelte/store';
import type { CorrectorsResponse } from '$lib/pocketbase-types';

/**
 * A writable Svelte store that holds a map of CorrectorsResponse records,
 * keyed by their ID. This store is not subscribed to PocketBase real-time updates.
 */
export const correctors = writable<Record<string, CorrectorsResponse>>({});
