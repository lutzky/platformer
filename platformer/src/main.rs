use bevy::prelude::*;
use bevy_debug_text_overlay::{screen_print, OverlayPlugin};
use bevy_pixel_camera::{PixelCameraPlugin, PixelViewport, PixelZoom};

fn main() {
    App::new()
        .add_plugins((
            DefaultPlugins.set(ImagePlugin::default_nearest()),
            OverlayPlugin {
                font_size: 23.0,
                ..default()
            },
        ))
        .add_plugins(PixelCameraPlugin)
        .add_systems(Startup, setup)
        // checkIsOnFloor
        // handleInput
        // moveX
        // handleXCollisions
        // moveY
        .add_systems(Update, handle_input_x)
        .add_systems(Update, handle_x_collisions)
        .add_systems(Update, handle_jump)
        .add_systems(Update, apply_gravity)
        .add_systems(Update, handle_y_collisions)
        .add_systems(Update, animate_sprite)
        .add_systems(Update, position_movement)
        .run();
}

const SUBPIXEL_RES: i16 = 128;
const SCREEN_WIDTH: i16 = 320;
const SCREEN_HEIGHT: i16 = 180;
const GRAVITY: i16 = 3;
const TERMINAL_VELOCITY: i16 = 10;

const DEFAULT_JUMP_SPEED: i16 = 2*SUBPIXEL_RES;
const DEFAULT_JUMP_HOVER_SPEED: i16 = SUBPIXEL_RES/2;

const TILE_MAP: [&str; 9] = [
    "..........",
    "..........",
    "..xxxx....",
    "..........",
    ".x.......x",
    "..x.....xx",
    "...x...xxx",
    "......xxxx",
    "xxxxxxxxxx",
];

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

    is_on_floor: bool,

    is_jumping: bool,
    jump_hover_speed: i16,
    jump_speed: i16,
    jump_started: bool,
}

impl Default for Player {
    fn default() -> Self {
        Self {
            jump_speed: DEFAULT_JUMP_SPEED,
            jump_hover_speed: DEFAULT_JUMP_HOVER_SPEED,
            vx: 0,
            vy: 0,

            is_on_floor: false,
            is_jumping: false,
            jump_started: false,
        }
    }
}

fn animate_sprite(time: Res<Time>, mut query: Query<(&mut AnimationTimer, &mut TextureAtlas)>) {
    for (mut timer, mut atlas) in query.iter_mut() {
        timer.timer.tick(time.delta());
        if timer.timer.just_finished() {
            atlas.index = (atlas.index + 1) % timer.frame_count;
        }
    }
}

fn handle_y_collisions(mut player: Query<(&mut Position, &mut Player)>) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    player.is_on_floor = false;

    // TODO: Check tiles

    if position.y <= -(SCREEN_HEIGHT / 2) * SUBPIXEL_RES + 16 * SUBPIXEL_RES {
        position.y = -(SCREEN_HEIGHT / 2) * SUBPIXEL_RES + 16 * SUBPIXEL_RES;
        player.vy = player.vy.max(0);
        player.is_on_floor = true;
    }

    if position.y >= ((SCREEN_HEIGHT / 2) - 16) * SUBPIXEL_RES {
        position.y = ((SCREEN_HEIGHT / 2) - 16) * SUBPIXEL_RES;
        player.vy = player.vy.min(0);
    }
}

fn handle_x_collisions(mut player: Query<(&mut Position, &mut Player)>) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    // TODO: Check tiles

    if position.x <= -(SCREEN_WIDTH / 2) * SUBPIXEL_RES + 16 * SUBPIXEL_RES {
        position.x = -(SCREEN_WIDTH / 2) * SUBPIXEL_RES + 16 * SUBPIXEL_RES;
        player.vx = player.vy.max(0);
        player.is_on_floor = true;
    }

    if position.x >= ((SCREEN_WIDTH / 2) - 16) * SUBPIXEL_RES {
        position.x = ((SCREEN_WIDTH / 2) - 16) * SUBPIXEL_RES;
        player.vx = player.vy.min(0);
    }
}

fn apply_gravity(time: Res<Time>, mut player: Query<(&mut Position, &mut Player)>) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    if player.is_on_floor {
        return;
    };

    player.vy = ((player.vy as f32 - (GRAVITY * SUBPIXEL_RES) as f32 * time.delta_seconds())
        as i16)
        .max(-TERMINAL_VELOCITY * SUBPIXEL_RES);
}

fn handle_jump(
    keyboard_input: Res<ButtonInput<KeyCode>>,
    mut player: Query<(&mut Position, &mut Player)>,
) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    if player.vy <= 0 {
        player.is_jumping = false;
    }

    if keyboard_input.pressed(KeyCode::Space) {
        if !player.jump_started {
            screen_print!(sec:0.5, "jump!");
            player.jump_started = true;
            if player.is_on_floor {
                player.vy = player.jump_speed;
                player.is_jumping = true;
            }
        }
    } else {
        player.jump_started = false;
        if player.is_jumping && player.vy > player.jump_hover_speed {
            player.vy = player.jump_hover_speed
        }
        player.is_jumping = false
    }
}

fn handle_input_x(
    keyboard_input: Res<ButtonInput<KeyCode>>,
    mut player: Query<(&mut Position, &mut Player)>,
) {
    let Ok((mut position, mut player)) = player.get_single_mut() else {
        return;
    };

    screen_print!(
        "player: ({:>5},{:>5}) v:({:>4},{:>4}) on_floor: {:>5}",
        position.x,
        position.y,
        player.vx,
        player.vy,
        player.is_on_floor
    );

    if keyboard_input.pressed(KeyCode::ArrowRight) {
        player.vx += 10;
    } else if keyboard_input.pressed(KeyCode::ArrowLeft) {
        player.vx -= 10;
    } else if player.vx > 0 {
        player.vx = (player.vx - 30).max(0);
    } else if player.vx < 0 {
        player.vx = (player.vx + 30).min(0);
    }

    position.x += player.vx;
    position.y += player.vy;
}

fn position_movement(mut query: Query<(&mut Transform, &Position)>) {
    for (mut transform, position) in query.iter_mut() {
        transform.translation.x = (position.x / SUBPIXEL_RES).into();
        transform.translation.y = (position.y / SUBPIXEL_RES).into();
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
            width: SCREEN_WIDTH.into(),
            height: SCREEN_HEIGHT.into(),
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
        Player { ..default() },
        Position { x: 0, y: 0 },
    ));
}
