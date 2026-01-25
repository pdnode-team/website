<script lang="ts">
	import { getLocale } from '$lib/paraglide/runtime.js';

	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let r: string | null = $state('');
  let disabledBtn = $state(false)

	onMount(() => {
		if (pb.authStore.isValid) {
			window.location.pathname = '/';
		}
		// 1. 获取查询字符串部分
		const queryString = window.location.search;

		// 2. 实例化 URLSearchParams
		const urlParams = new URLSearchParams(queryString);

		r = urlParams.get('r');
	});

	const handleLogin = async () => {
		error = '';
    disabledBtn = true
		// 处理登录逻辑
		try {
			const result = await pb.collection('users').authWithPassword(email.trim(), password.trim());
			console.log(result);
			if (r) {
				window.location.href = r;
			} else {
				window.location.pathname = '/';
			}
		} catch (e) {
      disabledBtn = false
			// @ts-ignore
			if (e.status === 400) {
				error = 'Check your email and password.';
			}
		}
	};
</script>

<svelte:head>
	<title>Login | Pdnode</title>
</svelte:head>

<div class="flex w-full h-full justify-center items-center flex-col gap-8 min-h-[80vh]">
	<div class="card w-full max-w-sm shadow-2xl bg-base-100 border border-base-200">
		<form
			class="card-body"
			onsubmit={(e) => {
				e.preventDefault;
				handleLogin();
			}}
		>
			<h1 class="text-3xl font-bold text-center mb-4">Login</h1>

			{#if getLocale() == 'zh'}
				<span class="alert alert-warning alert-soft flex justify-center items-center">
					<p class="text-center">此页面不提供中文版本</p>
				</span>
			{/if}
			{#if r && r.startsWith('/checkout')}
				<span class="alert alert-warning alert-soft flex justify-center items-center flex-col">
					<p class="text-center">After successful login, you will be redirected to checkout page</p>
				</span>
			{:else if r}
				<span class="alert alert-warning alert-soft flex justify-center items-center flex-col">
					<p class="text-center">After successful login, you will be redirected to:</p>
					<p class="text-ceter">{r}</p>
				</span>
			{/if}

			<div class="form-control">
				<label class="label" for="email">
					<span class="label-text font-medium">Email address</span>
				</label>
				<input
					id="email"
					type="email"
					placeholder="email@example.com"
					bind:value={email}
					class="input input-bordered focus:input-primary"
					required
				/>
			</div>

			<div class="form-control mt-4">
				<label class="label" for="password">
					<span class="label-text font-medium">Password</span>
				</label>
				<input
					id="password"
					type="password"
					placeholder="Enter your password"
					bind:value={password}
					class="input input-bordered focus:input-primary"
					required
				/>
			</div>

			{#if error}
				<span class="mt-2 alert alert-error alert-soft justify-center items-center flex">
					<p class="text-center">{error}</p>
				</span>
			{/if}

			<div class="form-control mt-6 flex justify-center items-center">
				<button type="submit" class="btn btn-primary no-animation text-white" disabled={disabledBtn}> Log in now </button>
			</div>

			<div class="text-center mt-4 text-sm">
				Don't have an account yet? <a href="/register" class="link link-primary font-semibold"
					>Go to register</a
				>
			</div>
		</form>
	</div>
</div>
