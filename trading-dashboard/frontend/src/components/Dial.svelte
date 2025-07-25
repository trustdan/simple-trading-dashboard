<script>
	export let value = 0; // Rating value (-3 to +3)
	export let size = 200; // Dial size in pixels
	export let label = ''; // Dial label
	export let interactive = true; // Whether the dial can be interacted with
	export let onValueChange = null; // Callback for value changes

	let dialElement;
	let isDragging = false;

	// Valid values (whole numbers from -3 to +3)
	const validValues = [-3, -2, -1, 0, 1, 2, 3];

	// Calculate dot position based on value (-3 to +3 maps to angles around the circle)
	$: dotAngle = (value / 3) * 135; // -135° to +135°
	$: dotPosition = getDotPosition(dotAngle);

	// Calculate color based on value (red to yellow to green)
	$: dotColor = getDotColor(value);

	function getDotPosition(angle) {
		const radians = (angle - 90) * (Math.PI / 180); // Convert to radians, -90° for top start
		const radius = size / 2 - 30; // Distance from center
		const x = size / 2 + Math.cos(radians) * radius;
		const y = size / 2 + Math.sin(radians) * radius;
		return { x, y };
	}

	function getDotColor(val) {
		if (val < 0) {
			// Red to orange for negative values
			const intensity = Math.abs(val) / 3;
			return `hsl(${20 - intensity * 20}, 80%, 60%)`; // Red to orange
		} else if (val === 0) {
			return '#666'; // Neutral gray
		} else {
			// Yellow to green for positive values
			const intensity = val / 3;
			return `hsl(${60 + intensity * 60}, 80%, 60%)`; // Yellow to green
		}
	}

	function handleMouseDown(event) {
		if (!interactive) return;
		isDragging = true;
		updateValueFromMouse(event);
	}

	function handleMouseMove(event) {
		if (!interactive || !isDragging) return;
		updateValueFromMouse(event);
	}

	function handleMouseUp() {
		isDragging = false;
	}

	function updateValueFromMouse(event) {
		if (!dialElement) return;

		const rect = dialElement.getBoundingClientRect();
		const centerX = rect.left + rect.width / 2;
		const centerY = rect.top + rect.height / 2;
		
		const x = event.clientX - centerX;
		const y = event.clientY - centerY;
		
		// Calculate angle (-180 to +180)
		let angle = Math.atan2(y, x) * (180 / Math.PI);
		
		// Adjust angle to our coordinate system (-135 to +135)
		angle = angle - 90; // Rotate so top is 0°
		
		// Clamp to our range
		if (angle > 135) angle = 135;
		if (angle < -135) angle = -135;
		
		// Convert angle to value and snap to nearest valid value
		const rawValue = (angle / 135) * 3;
		const newValue = validValues.reduce((prev, curr) => 
			Math.abs(curr - rawValue) < Math.abs(prev - rawValue) ? curr : prev
		);
		
		if (newValue !== value) {
			value = newValue;
			if (onValueChange) {
				onValueChange(value);
			}
		}
	}

	function handleKeyDown(event) {
		if (!interactive) return;

		switch (event.key) {
			case 'ArrowUp':
			case 'ArrowRight':
				event.preventDefault();
				const nextIndex = validValues.indexOf(value) + 1;
				if (nextIndex < validValues.length) {
					value = validValues[nextIndex];
					if (onValueChange) onValueChange(value);
				}
				break;
			case 'ArrowDown':
			case 'ArrowLeft':
				event.preventDefault();
				const prevIndex = validValues.indexOf(value) - 1;
				if (prevIndex >= 0) {
					value = validValues[prevIndex];
					if (onValueChange) onValueChange(value);
				}
				break;
		}
	}

	// Handle clicks on position markers to jump directly to that value
	function jumpToValue(targetValue) {
		if (!interactive || value === targetValue) return;
		value = targetValue;
		if (onValueChange) {
			onValueChange(value);
		}
	}
</script>

<svelte:window on:mousemove={handleMouseMove} on:mouseup={handleMouseUp} />

<div class="dial-container" style="width: {size}px; height: {size}px;">
	{#if label}
		<div class="dial-label">{label}</div>
	{/if}
	
	<svg
		bind:this={dialElement}
		width={size}
		height={size}
		class="dial-svg"
		class:interactive
		on:mousedown={handleMouseDown}
		on:keydown={handleKeyDown}
		tabindex={interactive ? 0 : -1}
	>
		<!-- Background circle -->
		<circle
			cx={size / 2}
			cy={size / 2}
			r={size / 2 - 40}
			fill="none"
			stroke="#333"
			stroke-width="4"
		/>

		<!-- Position markers for each valid value -->
		{#each validValues as markerValue}
			{@const markerAngle = (markerValue / 3) * 135}
			{@const markerPos = getDotPosition(markerAngle)}
			
			<circle
				cx={markerPos.x}
				cy={markerPos.y}
				r="6"
				fill={markerValue === value ? getDotColor(markerValue) : '#444'}
				stroke="#666"
				stroke-width="1"
				class="position-marker"
				class:interactive
				on:click={() => jumpToValue(markerValue)}
			/>
			
			<!-- Value labels -->
			{@const labelRadius = size / 2 - 15}
			{@const labelRadians = (markerAngle - 90) * (Math.PI / 180)}
			{@const labelX = size / 2 + Math.cos(labelRadians) * labelRadius}
			{@const labelY = size / 2 + Math.sin(labelRadians) * labelRadius}
			
			<text
				x={labelX}
				y={labelY}
				text-anchor="middle"
				dominant-baseline="middle"
				fill="#888"
				font-size="12"
				font-weight="500"
				class="value-label"
			>
				{markerValue > 0 ? '+' : ''}{markerValue}
			</text>
		{/each}

		<!-- Moving dot -->
		<circle
			cx={dotPosition.x}
			cy={dotPosition.y}
			r="12"
			fill={dotColor}
			stroke="#fff"
			stroke-width="2"
			class="dot"
			class:dragging={isDragging}
		/>

		<!-- Center value display -->
		<text
			x={size / 2}
			y={size / 2}
			text-anchor="middle"
			dominant-baseline="middle"
			fill={dotColor}
			font-size="20"
			font-weight="bold"
			class="center-value"
		>
			{value > 0 ? '+' : ''}{value}
		</text>
	</svg>
</div>

<style>
	.dial-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
	}

	.dial-label {
		font-size: 14px;
		font-weight: 600;
		color: #ccc;
		text-align: center;
		margin-bottom: 4px;
	}

	.dial-svg {
		transition: all 0.2s ease;
	}

	.dial-svg.interactive {
		cursor: grab;
	}

	.dial-svg.interactive:active {
		cursor: grabbing;
	}

	.dial-svg.interactive:focus {
		outline: 2px solid #4a90e2;
		outline-offset: 4px;
	}

	.dot {
		transition: all 0.2s ease-in-out;
		filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
	}

	.dot.dragging {
		filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.4));
		transform-origin: center;
	}

	.position-marker.interactive {
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.position-marker.interactive:hover {
		r: 8;
		stroke-width: 2;
	}

	.value-label {
		font-family: system-ui, -apple-system, sans-serif;
		pointer-events: none;
	}

	.center-value {
		font-family: system-ui, -apple-system, sans-serif;
		pointer-events: none;
	}
</style> 