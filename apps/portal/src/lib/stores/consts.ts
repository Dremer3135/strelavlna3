import { writable } from 'svelte/store';
import type { EditableConst } from '$lib/types';

/**
 * A writable Svelte store that holds a map of problem records, keyed by the problem ID.
 * Its initial value is an empty object.
 */
export const editableConsts = writable<Record<string, EditableConst>>({});
