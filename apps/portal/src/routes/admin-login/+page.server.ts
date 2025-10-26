import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals, url }) => {
        const formData = await request.formData();
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;
        const redirectTo = formData.get('redirectTo');

        try {
            await locals.pb.collection('correctors').authWithPassword(email, password);
        } catch (err) {
            console.error(err);
            return fail(400, { error: 'Invalid email or password.' });
        }

        if (redirectTo) {
            throw redirect(303, redirectTo.toString());
        }

        throw redirect(303, '/admin');
    }
};

// We need to combine the load function and the actions into the same file.
export const load = ({ url, locals }) => {
    if (locals.user && locals.user.collectionName === 'correctors') {
        throw redirect(303, '/admin');
    }
    
    const redirectTo = url.searchParams.get('redirectTo');
    return { redirectTo };
};
