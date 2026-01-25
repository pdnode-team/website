<script lang="ts">
	import { onDestroy, onMount } from 'svelte';

	let status = $state('idle');

	let modalRef: HTMLDialogElement | null = $state(null);

  let intervalId: ReturnType<typeof setInterval>; 

	async function checkStatus() {
    try {
      const response = await fetch('http://127.0.0.1:8090/api/health');

      if (response.ok) {
        status = 'healthy';
        // 只有当弹窗开着的时候才去关闭
        if (modalRef?.open) {
          modalRef.close();
					window.location.reload()
        }
      } else {
        throw new Error('Unhealthy');
      }
    } catch (err) {
      status = 'unhealthy';
      // 关键：只有当弹窗没打开时，才尝试打开它
      if (modalRef && !modalRef.open) {
        modalRef.showModal();
      }
    }
  }

	onMount(async () => {
		checkStatus();

		intervalId = setInterval(() => {
      // 如果用户正在手动点击重试（status === 'checking'），可以跳过这次自动检查
      if (status !== 'checking') {
        checkStatus();
      }
    }, 5000);
	});

  onDestroy(() => {
    // 组件卸载时必须清除定时器，否则后台会一直请求
    if (intervalId) clearInterval(intervalId);
  });

	async function handleCheckStatus(e: Event) {
		e.preventDefault();
		status = 'checking';
		await new Promise((resolve) => setTimeout(resolve, 400));
		checkStatus();
	}

	function handleCancel(e: Event) {
		// 核心：阻止 ESC 键关闭弹窗
		e.preventDefault();
	}
</script>

<dialog bind:this={modalRef} class="modal" oncancel={(e) => handleCancel(e)}>
	<div class="modal-box border-2 border-error">
		<h3 class="text-lg font-bold">The Backend API is unhealthy</h3>
		<p class="py-4">
			We are unable to connect to the backend API, or the backend API is not ready to accept
			requests.
		</p>
		<p>
			If you don't need to log in, register, view announcements, or subscribe, <a
				href="/"
				class="link">please click to return to the homepage.</a
			>
		</p>

		<div class="modal-action">
			<button class="btn" onclick={(e) => handleCheckStatus(e)} disabled={status === 'checking'}>
				Recheck
			</button>
		</div>
	</div>
</dialog>
