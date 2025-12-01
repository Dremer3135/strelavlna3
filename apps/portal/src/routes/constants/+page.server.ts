import type { ConstantsResponse } from '$lib/pocketbase-types';

export async function load({locals}) {
  let constants: ConstantsResponse[] = await locals.pb.collection('constants').getFullList();

  return {
    constants: constants
  }
}