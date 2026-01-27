<script lang="ts">
  import { getLocale } from '$lib/paraglide/runtime.js';

  import pb from '$lib/pocketbase';
  import { onMount } from 'svelte';

  let email = $state('');
  let password = $state('');
  let name = $state('');

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
    if (r) {
      r = r.replace(/(^\w+:|^)\/\//, '');
    }
  });

  const handleRegister = async () => {
    error = '';
    disabledBtn = true

    try {
      const result = await pb.collection('users').create({
        email: email.trim(),
        password: password.trim(),
        passwordConfirm: password.trim(),
        name: name.trim(),
      });
      console.log(result);
      if (r) {
        window.location.href = r;
      } else {
        window.location.pathname = '/login';
      }
    } catch (e) {
      disabledBtn = false
      // @ts-ignore
      if (e.status === 400) {
        // @ts-ignore
        switch (e.data?.data?.email?.code) {
          case "validation_not_unique":
            error = "The email address is already in use."
            break
          case "validation_is_email":
            error = "Must be a valid email address."
            break
          default:
            console.log(e)

            error = "Please check your input.";
            break
        }
        // @ts-ignore
        switch (e.data?.data?.password?.code) {
          case "validation_min_text_constraint":
            error = "Password must be at least 8 characters."
            break
          default:
            console.log(e)

            error = "Please check your input.";
            break
        }
      }
      else {
        error = "unknown error(s)"
      }
    }
  };
</script>

<svelte:head>
  <title>Register | Pdnode</title>
</svelte:head>

<div class="flex w-full h-full justify-center items-center flex-col gap-8 min-h-[80vh]">
  <div class="card w-full max-w-sm shadow-2xl bg-base-100 border border-base-200">
    <form
      class="card-body"
      onsubmit={(e) => {
				e.preventDefault;
				handleRegister();
			}}
    >
      <h1 class="text-3xl font-bold text-center mb-4">Register</h1>

      {#if getLocale() === 'zh'}
				<span class="alert alert-warning alert-soft flex justify-center items-center">
					<p class="text-center">此页面不提供中文版本</p>
				</span>
      {/if}
      {#if r && r.startsWith('/checkout')}
				<span class="alert alert-warning alert-soft flex justify-center items-center flex-col">
					<p class="text-center">After successful register, you will be redirected to checkout page</p>
				</span>
      {:else if r}
				<span class="alert alert-warning alert-soft flex justify-center items-center flex-col">
					<p class="text-center">After successful register, you will be redirected to:</p>
					<p class="text-ceter">{r}</p>
				</span>
      {/if}

      <div class="form-control">
        <label class="label" for="name">
          <span class="label-text font-medium">Name</span>
        </label>
        <input
          id="name"
          type="text"
          placeholder="John Doe"
          bind:value={name}
          class="input input-bordered focus:input-primary"
          required
        />
      </div>

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
        <button type="submit" class="btn btn-primary no-animation text-white" disabled={disabledBtn}> Register now </button>
      </div>

      <div class="text-center mt-4 text-sm">
        Have an account? <a href="/login" class="link link-primary font-semibold"
      >Go to Login</a
      >
      </div>
    </form>
  </div>
</div>
