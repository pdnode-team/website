<script lang="ts">
  import {onMount} from "svelte";
  import pb from "$lib/pocketbase";

  let subscription: any = $state()

  const getHash = (str: string) => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    return hash;
  };

  // 生成 HSL 颜色字符串
  // 360 是色相环的全周
  const generateColor = (str: string) => {
    if (!str) return `background-color: #3d4451`;
    const hash = getHash(str);
    const hue = Math.abs(hash % 360);
    return `background-color: hsl(${hue}, 70%, 45%)`;
  };

  onMount(async () => {
    if (!pb.authStore.isValid) {
      window.location.pathname = '/login';
    }

    const subscribe = await fetch("/api/subscription", {
      headers: { 'content-type': 'application/json', Authorization: pb.authStore.token }
    })

    const json = await subscribe.json()

    if (subscribe.status == 400){
      subscription = null
    }else if (subscribe.ok){
      subscription = json
    }

    console.log(json)

  })
</script>

<svelte:head>
  <title>Profile | Pdnode</title>
</svelte:head>

<div class="flex w-full h-full justify-center items-center flex-col gap-8 min-h-[80vh] p-4">
  <div class="card w-full max-w-sm shadow-2xl bg-base-100 border border-base-200">
    <div class="card-body items-center text-center">
      <h1 class="card-title text-3xl font-bold mb-6">Profile</h1>

      <div class="avatar placeholder mb-4 mt-4">
        <div class="bg-neutral text-neutral-content rounded-full w-24 flex items-center justify-center" style={generateColor(pb.authStore.record?.name)}>
          <span class="text-3xl">{pb.authStore.record?.name[0]}</span>
        </div>
      </div>

      <div class="w-full space-y-4">
        <div class="form-control">
            <span class="label-text font-semibold">Name</span>
          <div class="bg-base-200 p-3 rounded-lg text-left flex items-center justify-start gap-4">
            {pb.authStore.record?.name}
            {#if subscription}
              <div class="badge badge-primary">{subscription.plan}</div>
            {/if}


          </div>
        </div>

        <div class="form-control">
            <span class="label-text font-semibold">Email</span>
          <div class="bg-base-200 p-3 rounded-lg text-left">
            {pb.authStore.record?.email}
          </div>
        </div>
      </div>

      <div class="card-actions justify-end mt-6 w-full">
        <button class="btn btn-primary btn-block" onclick={() => {pb.authStore.clear(); window.location.reload()}}>Log Out</button>

      </div>
    </div>
  </div>
</div>