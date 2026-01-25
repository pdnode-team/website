<script lang="ts">
	import pb from '$lib/pocketbase';
	import { PUBLIC_POCKETABSE_URL } from '$env/static/public';
	import { onMount } from 'svelte';

	let errorMessage = $state('');

	onMount(() => {
		if (!pb.authStore.isValid) {
			window.location.href = `/login?r=${window.location.pathname}${window.location.search}`;
		}

		// 1. 获取查询字符串部分
		const queryString = window.location.search;

		// 2. 实例化 URLSearchParams
		const urlParams = new URLSearchParams(queryString);

		plan = urlParams.get('plan');

		handleCheckout();
	});

	let plan: string | null = $state('');

	async function handleCheckout() {
		const response = await fetch(`${PUBLIC_POCKETABSE_URL}/api/checkout/subscription`, {
			method: 'POST',
			body: JSON.stringify({ plan }),
			headers: { 'content-type': 'application/json', Authorization: pb.authStore.token }
		});

		const { url, message, status } = await response.json();
		if (url) {
			window.location.href = url; // 跳转到 Stripe 支付页面
		} else {
			if (status === 400) {
				errorMessage = message;
			} else if (status === 500) {
				errorMessage = 'Internal server error.';
        console.error(message)
			}
		}
	}
</script>

<div class="flex justify-center items-center w-full h-screen">
	<span class="flex justify-center items-center alert alert-soft alert-error flex-col">
		<p class="text-center">
			{errorMessage}
		</p>
	</span>
</div>
