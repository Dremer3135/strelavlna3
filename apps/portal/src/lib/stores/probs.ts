import { writable } from 'svelte/store';
import type { EditableProb } from '$lib/types';
import { pocketbase } from '$lib/pocketbase';
import type { ProbsResponse } from '$lib/types/pocketbase-types';

/**
 * A writable Svelte store that holds a map of problem records, keyed by the problem ID.
 * Its initial value is an empty object.
 */
export const editableProbs = writable<Record<string, EditableProb>>({});

export async function subscribeToProbs(initialData: ProbsResponse[]) {
	const initialProbs = Object.fromEntries(
		initialData.map((prob) => [prob.id, { prob: prob, edit: {} }])
	);
	editableProbs.set(initialProbs);

	console.log('Subscribing to probs with auth model:', pocketbase.authStore.model);
	const unsubscribe = await pocketbase.collection('probs').subscribe('*', (e) => {
		const record = e.record as ProbsResponse;
		if (e.action === 'create') {
			editableProbs.update((currentProbs) => {
				currentProbs[record.id] = { prob: record, edit: {} };
				return currentProbs;
			});
		} else if (e.action === 'update') {
			editableProbs.update((currentProbs) => {
				if (currentProbs[record.id]) {
					const updatedEditableProb = {
						...currentProbs[record.id],
						prob: record
					};
					return {
						...currentProbs,
						[record.id]: updatedEditableProb
					};
				}
				return currentProbs;
			});
		} else if (e.action === 'delete') {
			editableProbs.update((currentProbs) => {
				delete currentProbs[record.id];
				return currentProbs;
			});
		}
	});

	return unsubscribe;
}
