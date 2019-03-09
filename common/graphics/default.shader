#shader vertex
#version 130

in vec3 vertex_pos;
in vec2 vertex_uv;
in vec3 vertex_normal;

out vec2 uv;
out vec3 normal;

struct LightSoruce {
	int Type;
	vec3 Position;
	vec3 Direction;
	vec3 Color;
	float Attenuation;
};

uniform mat4 u_MVP;
//uniform LightSource u_Lights[16];

void main() {
	uv = vertex_uv;
	normal = vertex_normal;
    gl_Position = u_MVP * vec4(vertex_pos, 1.0);
}

#shader fragment
#version 130

in vec2 uv;
in vec3 normal;

out vec4 fColor;

uniform sampler2D tex;

void main() {
    fColor = texture(tex, uv);
}
