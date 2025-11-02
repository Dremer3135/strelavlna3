import { fail, redirect, error } from '@sveltejs/kit';
import { enhance } from '$app/forms';
import { createPocketbaseInstance } from '$lib/server/pocketbase';

import { POCKETBASE_SUPERUSER_EMAIL, POCKETBASE_SUPERUSER_PASSWORD } from '$env/static/private';


export const load = async ({ locals }) => {
	try {
		const probs = await locals.pb.collection('probs').getFullList();
		const consts = await locals.pb.collection('consts').getFullList();

		return {
			probs: probs,
			consts: consts
            // pb: locals.pb
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}
};
