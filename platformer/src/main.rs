use bevy::prelude::*;
use bevy::window::WindowResolution;

fn main() {
    App::new()
        .add_plugins(DefaultPlugins.set(ImagePlugin::default_nearest()).set(WindowPlugin {
            primary_window: Some(Window {
                resolution: WindowResolution::new(320.,200.).with_scale_factor_override(6.0),
                ..default()
            }),
            ..default()
        }))
        .add_systems(Startup, setup)
        .add_systems(Update, animate_sprite)
        .add_systems(Update, player_movement)
        .run();
}

#[derive(Component)]
struct AnimationTimer {
    timer: Timer,
    frame_count: usize
}

#[derive(Component)]
struct Player;

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
    mut player: Query<&mut Transform, With<Player>>,
) {
    let Ok(mut transform) = player.get_single_mut() else {
        return;
    };

    if keyboard_input.pressed(KeyCode::ArrowUp) {
        transform.translation += Vec3::new(0.0, 1.0, 0.0);
    }
    if keyboard_input.pressed(KeyCode::ArrowDown) {
        transform.translation += Vec3::new(0.0,- 1.0, 0.0);
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

    commands.spawn(Camera2dBundle::default());
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
        Player {},
    ));
}
