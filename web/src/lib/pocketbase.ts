import PocketBase from 'pocketbase';
import { PUBLIC_POCKETABSE_URL } from '$env/static/public';

const pb = new PocketBase(PUBLIC_POCKETABSE_URL);

export default pb;