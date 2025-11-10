import { writable } from 'svelte/store';
import type { EditableConstant } from '$lib/types';
import { pocketbase } from '$lib/pocketbase';
import type { ConstantsResponse } from '$lib/types/pocketbase-types';

/**
 * A writable Svelte store that holds a map of problem records, keyed by the problem ID.
 * Its initial value is an empty object.
 */
export const editableConstants = writable<Record<string, EditableConstant>>({});

export async function subscribeToConstants(initialData: ConstantsResponse[]) {
	const initialConsts = Object.fromEntries(
		initialData.map((c) => [c.id, { constant: c, edit: {} }])
	);
	editableConstants.set(initialConsts);

	const unsubscribe = await pocketbase.collection('constants').subscribe('*', (e) => {
		const record = e.record as ConstantsResponse;
		if (e.action === 'create') {
			editableConstants.update((currentConsts) => {
				currentConsts[record.id] = { constant: record, edit: {} };
				return currentConsts;
			});
		} else if (e.action === 'update') {
			editableConstants.update((currentConsts) => {
				if (currentConsts[record.id]) {
					const updatedEditableConst = {
						...currentConsts[record.id],
						constant: record
					};
					return {
						...currentConsts,
						[record.id]: updatedEditableConst
					};
				}
				return currentConsts;
			});
		} else if (e.action === 'delete') {
			editableConstants.update((currentConsts) => {
				delete currentConsts[record.id];
				return currentConsts;
			});
		}
	});

	return unsubscribe;
}
