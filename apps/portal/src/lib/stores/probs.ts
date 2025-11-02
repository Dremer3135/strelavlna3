import { writable } from 'svelte/store';
import type { EditableProb } from '$lib/types';

/**
 * A writable Svelte store that holds a map of problem records, keyed by the problem ID.
 * Its initial value is an empty object.
 */
export const editableProbs = writable<Record<string, EditableProb>>({});
