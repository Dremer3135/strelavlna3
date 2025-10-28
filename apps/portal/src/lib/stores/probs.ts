import { writable } from 'svelte/store';
import type { ProbsResponse } from '../pocketbase-types';

/**
 * A writable Svelte store that holds a list of problem records.
 * Its initial value is an empty array.
 */
export const probs = writable<ProbsResponse[]>([]);
