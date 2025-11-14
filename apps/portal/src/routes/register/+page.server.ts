import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals, url }) => {
        const formData = await request.formData();
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;
        const passwordConfirm = formData.get('password-confirm') as string;
        const school = formData.get('school');
        const redirectTo = formData.get('redirectTo') as string;

        try {
            
            await locals.pb.collection('teachers').create({
                email,
                password,
                passwordConfirm: passwordConfirm,
                skola: school,
            });
        } catch (err) {
            console.error(err);
            return fail(400, { error: 'Invalid email or password.' });
        }

        if (redirectTo) {
            throw redirect(303, redirectTo);
        }
    }
};

// We need to combine the load function and the actions into the same file.
export const load = ({ url, locals }) => {
    if (locals.user && (locals.user.collectionName === 'correctors' || locals.user.collectionName === 'teachers')) {
        throw redirect(303, '/');
    }
    const redirectTo = url.searchParams.get('redirectTo');
    return { redirectTo };
};
