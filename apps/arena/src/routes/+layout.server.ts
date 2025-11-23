import type { ContestState } from "$lib/types";
import { redirect } from "@sveltejs/kit";

export const load = async ({ params, url }) => {
    let state: ContestState = "waiting";

    if (state == "waiting") {
        throw redirect(303, "/waitroom");
    } else if (state == "ended") {
        throw redirect(303, "/results");
    }

    let token = url.searchParams.get('token');
    return { token };
}
