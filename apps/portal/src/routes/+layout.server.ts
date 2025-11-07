export const load = async ({ locals }) => {
	return {
		user: locals.user,
		cookie: locals.pb.authStore.exportToCookie()
	};
};