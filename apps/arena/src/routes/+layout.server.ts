import type { ContestState } from "$lib/types";
import { redirect } from "@sveltejs/kit";

export const load = async () => {
    let state: ContestState = "waiting";

    if (state == "waiting") {
        throw redirect(303, "/waitroom");
    } else if (state == "ended") {
        throw redirect(303, "/results");
    }
}