import { writable } from 'svelte/store';
import type { Prob } from '$lib/types';

export const probs = writable<Record<string, Prob>>({});
