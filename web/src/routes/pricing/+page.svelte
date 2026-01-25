<div class="flex justify-center items-center h-screen w-full flex-col gap-5">
	<h1 class="text-5xl font-bold">Pricing</h1>
	<p class="text-2xl">Subscribing to us is the best way to support Pdnode!</p>
	<div class="flex justify-center items-center gap-4 max-md:flex-col">
		<div class="card w-96 bg-base-100 shadow-lg">
			<div class="card-body">
				<div class="flex justify-between">
					<h2 class="text-3xl font-bold">Starter</h2>
					<span class="text-xl">$5/mo</span>
				</div>
				<ul class="mt-6 flex flex-col gap-2 text-xs">
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Multiple services in the future</span>
					</li>
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Faster support services</span>
					</li>
				</ul>
				<div class="mt-6">
					<button
						class="btn btn-primary btn-block"
						onclick={() => (window.location.href = '/checkout?plan=starter')}>Subscribe</button
					>
				</div>
			</div>
		</div>
		<div class="card w-96 bg-base-100 shadow-sm">
			<div class="card-body">
				<!-- <span class="badge badge-xs badge-warning">Most Popular</span> -->
				<div class="flex justify-between">
					<h2 class="text-3xl font-bold">Plus</h2>
					<span class="text-xl">$15/mo</span>
				</div>
				<ul class="mt-6 flex flex-col gap-2 text-xs">
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Multiple services in the future</span>
					</li>
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Faster support services</span>
					</li>
				</ul>
				<div class="mt-6">
					<button
						class="btn btn-primary btn-block"
						onclick={() => (window.location.href = '/checkout?plan=plus')}>Subscribe</button
					>
				</div>
			</div>
		</div>
		<div class="card w-96 bg-base-100 shadow-sm">
			<div class="card-body">
				<div class="flex justify-between">
					<h2 class="text-3xl font-bold">Pro</h2>
					<span class="text-xl">$25/mo</span>
				</div>
				<ul class="mt-6 flex flex-col gap-2 text-xs">
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Multiple services in the future</span>
					</li>
					<li>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="size-4 me-2 inline-block text-success"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							/></svg
						>
						<span>Faster support services</span>
					</li>
				</ul>
				<div class="mt-6">
					<button
						class="btn btn-primary btn-block"
						onclick={() => (window.location.href = '/checkout?plan=pro')}>Subscribe</button
					>
				</div>
			</div>
		</div>
	</div>
	<p class="text-sm">*All prices exclude tax.</p>
</div>
