<script>
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

	onMount(() => {
		if (!pb.authStore.isValid) {
			window.location.href = `/login?r=${window.location.pathname}${window.location.search}`;
		}
	});
</script>

<div class="flex justify-center items-center w-full h-screen">
	<span class="flex justify-center items-center alert alert-soft alert-success flex-col">
		{#if pb.authStore.record?.name}
			<h3 class="text-2xl font-bold">Congratulations {pb.authStore.record?.name}!</h3>
		{:else}
			<h3 class="text-2xl font-bold">Congratulations!</h3>
		{/if}

		<p class="text-center">
			Your subscription has been successfully completed! Your profile will be update within one
			minute.
		</p>
		<p>Subscriptions may be delayed by up to one minute.</p>
		<p class="font-bold">You can now return to the homepage or close this window!</p>
	</span>
</div>
