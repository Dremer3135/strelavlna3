import PocketBase from 'pocketbase';
import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import type { TypedPocketBase } from './pocketbase-types';

export const pb = new PocketBase(PUBLIC_POCKETBASE_URL) as TypedPocketBase;
