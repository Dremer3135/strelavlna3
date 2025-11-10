import { writable } from 'svelte/store';
import type { EditableContest } from '$lib/types';
import { pocketbase } from '$lib/pocketbase';
import type { ContestsResponse } from '$lib/pocketbase-types';

/**
 * A writable Svelte store that holds a map of contest records, keyed by the contest ID.
 * Its initial value is an empty object.
 */
export const editableContests = writable<Record<string, EditableContest>>({});

export async function subscribeToContests(initialData: ContestsResponse[]) {
	const initialContests = Object.fromEntries(
		initialData.map((contest) => [contest.id, { contest: contest, edit: {} }])
	);
	editableContests.set(initialContests);

	const unsubscribe = await pocketbase.collection('contests').subscribe('*', (e) => {
		const record = e.record as ContestsResponse;
		if (e.action === 'create') {
			editableContests.update((currentContests) => {
				currentContests[record.id] = { contest: record, edit: {} };
				return { ...currentContests }; // Return a new object to ensure reactivity
			});
		} else if (e.action === 'update') {
			editableContests.update((currentContests) => {
				if (currentContests[record.id]) {
					const updatedEditableContest = {
						...currentContests[record.id],
						contest: record
					};
					return {
						...currentContests,
						[record.id]: updatedEditableContest
					};
				}
				return currentContests;
			});
		} else if (e.action === 'delete') {
			editableContests.update((currentContests) => {
				delete currentContests[record.id];
				return { ...currentContests }; // Return a new object
			});
		}
	});

	return unsubscribe;
}
