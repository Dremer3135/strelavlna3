import { fail, redirect, error } from '@sveltejs/kit';
import { enhance } from '$app/forms';
import { createPocketbaseInstance } from '$lib/server/pocketbase';

import { POCKETBASE_SUPERUSER_EMAIL, POCKETBASE_SUPERUSER_PASSWORD } from '$env/static/private';


export const load = async ({ locals }) => {
	console.log("loading with: ", locals.pb.authStore);
	
	try {
		console.log("- - - - - - - - - - - - - - - - - - - - - - -");
		const probs = await locals.pb.collection('probs').getFullList();
		console.log("heheheha");
		const constants = await locals.pb.collection('constants').getFullList();
		console.log("lalalalalalala");
		
		// console.log(constants);
		return {
			probs: probs,
			constants: constants
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}
};
