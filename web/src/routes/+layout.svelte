<script lang="ts">
	import './layout.css';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let { children } = $props();

</script>

<svelte:head>
<title>Pdnode</title>
</svelte:head>
<div class="flex flex-col min-h-screen min-w-full">
  
  <Header />
  
  <main class="grow">
    {@render children()}
  </main>

  <Footer />
  
</div>