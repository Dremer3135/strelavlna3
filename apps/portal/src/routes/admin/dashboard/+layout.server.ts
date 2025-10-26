import { redirect } from '@sveltejs/kit'
import { isCorrector } from '$lib/utils';

export const load = async ({ locals }) => {
    // Use the type guard to check the user type.
    // If the user is a corrector and has the admin flag, allow access.
    // Otherwise, redirect.
    if (isCorrector(locals.user) && locals.user.admin === true) {
        // User is a true admin, do nothing and let the page load.
    } else {
        throw redirect(303, "/admin");
    }
}   