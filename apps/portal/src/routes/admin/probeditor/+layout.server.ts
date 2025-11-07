import { fail, redirect, error } from '@sveltejs/kit';
import { enhance } from '$app/forms';
import { createPocketbaseInstance } from '$lib/server/pocketbase';

import { POCKETBASE_SUPERUSER_EMAIL, POCKETBASE_SUPERUSER_PASSWORD } from '$env/static/private';


export const load = async ({ locals }) => {
	try {
		const probs = await locals.pb.collection('probs').getFullList();
		const constants = await locals.pb.collection('constants').getFullList();

		return {
			probs: probs,
			constants: constants
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}
};
