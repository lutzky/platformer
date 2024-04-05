use bevy::prelude::*;
use bevy::window::WindowResolution;
use bevy_pixel_camera::{PixelCameraPlugin, PixelViewport, PixelZoom};

fn main() {
    App::new()
        .add_plugins(DefaultPlugins.set(ImagePlugin::default_nearest()))
        .add_plugins(PixelCameraPlugin)
        .add_systems(Startup, setup)
        .add_systems(Update, animate_sprite)
        .add_systems(Update, player_movement)
        .add_systems(Update, position_movement)
        .run();
}

#[derive(Component)]
struct AnimationTimer {
    timer: Timer,
    frame_count: usize,
}

#[derive(Component)]
struct Position {
    x: i16,
    y: i16,
}

#[derive(Component)]
struct Player {
    vx: i16,
    vy: i16,
}

fn animate_sprite(time: Res<Time>, mut query: Query<(&mut AnimationTimer, &mut TextureAtlas)>) {
    for (mut timer, mut atlas) in query.iter_mut() {
        timer.timer.tick(time.delta());
        if timer.timer.just_finished() {
            atlas.index = (atlas.index + 1) % timer.frame_count;
        }
    }
}

fn player_movement(
    keyboard_input: Res<ButtonInput<KeyCode>>,
    mut player: Query<(&mut Position, &mut Player)>,
) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    if keyboard_input.pressed(KeyCode::ArrowRight) {
        player.vx += 10;
    } else if keyboard_input.pressed(KeyCode::ArrowLeft) {
        player.vx -= 10;
    } else if player.vx > 0 {
        player.vx = (player.vx - 30).max(0);
    } else if player.vx < 0 {
        player.vx = (player.vx + 30).min(0);
    }

    player.vy = player.vy.clamp(-512, 512);
    player.vx = player.vx.clamp(-512, 512);

    position.x += player.vx;
    position.y += player.vy;

    if !(-100 * 160..100 * 160).contains(&position.x) {
        player.vx = 0;
    }
    if !(-100 * 90..100 * 90).contains(&position.y) {
        player.vy = 0;
    }
    position.x = position.x.clamp(-100 * 160, 100 * 160);
    position.y = position.y.clamp(-100 * 90, 100 * 90);
}

fn position_movement(mut query: Query<(&mut Transform, &Position)>) {
    for (mut transform, position) in query.iter_mut() {
        transform.translation.x = (position.x / 256).into();
        transform.translation.y = (position.y / 256).into();
    }
}

fn setup(
    mut commands: Commands,
    asset_server: Res<AssetServer>,
    mut texture_atlas_layouts: ResMut<Assets<TextureAtlasLayout>>,
) {
    let layout = TextureAtlasLayout::from_grid(Vec2::new(32.0, 32.0), 11, 1, None, None);
    let texture = asset_server.load("sprites/Idle (32x32).png");
    let texture_atlas_layout = texture_atlas_layouts.add(layout);

    commands.spawn((
        Camera2dBundle::default(),
        PixelZoom::FitSize {
            width: 320,
            height: 180,
        },
        PixelViewport,
    ));
    commands.spawn((
        SpriteSheetBundle {
            texture,
            atlas: TextureAtlas {
                layout: texture_atlas_layout,
                index: 0,
            },
            // transform: Transform::from_scale(Vec3::splat(1.0)),
            ..default()
        },
        AnimationTimer {
            timer: Timer::from_seconds(0.1, TimerMode::Repeating),
            frame_count: 11,
        },
        Player { vx: 0, vy: 0 },
        Position { x: 0, y: 0 },
    ));
}
