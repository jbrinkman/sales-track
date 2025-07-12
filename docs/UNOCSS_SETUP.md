# UnoCSS Configuration for Sales Track

This document describes the UnoCSS setup and configuration for the Sales Track application.

## Overview

UnoCSS is configured as the primary CSS framework for the Sales Track application, providing:
- Utility-first CSS classes
- Custom theme colors and design tokens
- Responsive design utilities
- Custom shortcuts for common patterns

## Configuration Files

### 1. `frontend/uno.config.ts`
Main UnoCSS configuration file containing:
- **Presets**: UnoCSS core, Attributify, and Icons
- **Theme**: Custom color palette, typography, spacing
- **Shortcuts**: Reusable component patterns

### 2. `frontend/vite.config.ts`
Vite configuration with UnoCSS plugin integration.

### 3. `frontend/src/main.ts`
UnoCSS virtual import: `import 'virtual:uno.css'`

## Custom Theme

### Color Palette
- **Primary**: Blue color scheme (#3b82f6 family)
- **Success**: Green color scheme (#22c55e family)
- **Warning**: Yellow/Orange color scheme (#f59e0b family)
- **Error**: Red color scheme (#ef4444 family)

### Typography
- **Font Family**: Nunito (with system fallbacks)
- **Font Sizes**: Extended scale from xs to 4xl
- **Line Heights**: Optimized for readability

### Custom Shortcuts
- `btn`: Base button styling
- `btn-primary`: Primary button with hover/focus states
- `btn-secondary`: Secondary button styling
- `card`: Card component with shadow and border
- `input`: Form input styling with focus states

## Usage Examples

### Buttons
```vue
<button class="btn-primary">Primary Action</button>
<button class="btn-secondary">Secondary Action</button>
<button class="btn bg-success-600 text-white hover:bg-success-700">Success</button>
```

### Cards
```vue
<div class="card p-6">
  <h3 class="text-lg font-semibold mb-2">Card Title</h3>
  <p class="text-gray-600">Card content</p>
</div>
```

### Form Inputs
```vue
<input type="text" class="input w-full" placeholder="Enter text">
```

### Layout
```vue
<div class="container mx-auto px-4">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <!-- Responsive grid content -->
  </div>
</div>
```

## Testing

The UnoCSS configuration includes a test component (`UnoTest.vue`) that demonstrates:
- Button variations
- Card styling
- Input fields
- Color palette
- Responsive design

## Build Integration

UnoCSS is fully integrated with the Wails build process:
1. Vite processes UnoCSS during frontend compilation
2. Generated CSS is included in the final application bundle
3. No additional build steps required

## Best Practices

1. **Use Utility Classes**: Prefer utility classes over custom CSS
2. **Leverage Shortcuts**: Use custom shortcuts for repeated patterns
3. **Follow Theme**: Use theme colors and spacing for consistency
4. **Responsive Design**: Use responsive prefixes (sm:, md:, lg:, xl:)
5. **Component Scoping**: Use scoped styles only when necessary

## Troubleshooting

### Build Issues
- Ensure `virtual:uno.css` is imported in `main.ts`
- Check that UnoCSS plugin is properly configured in `vite.config.ts`
- Verify all UnoCSS dependencies are installed

### Style Not Applied
- Check class names for typos
- Ensure classes are included in the UnoCSS scan
- Use browser dev tools to verify CSS generation

## Future Enhancements

- Add icon collections for UnoCSS Icons preset
- Implement dark theme support
- Add animation utilities
- Create additional component shortcuts
