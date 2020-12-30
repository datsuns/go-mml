/* C example that opens a game music file and records 10 seconds to "out.wav" */

#include "gme/gme.h"

#include "Wave_Writer.h" /* wave_ functions for writing sound file */
#include <stdlib.h>
#include <stdio.h>

int handle_error( const char* str );

int nfs2wav(char* src, char* dest)
{
	long sample_rate = 44100; /* number of samples per second */
	int track = 0; /* index of track to play (0 = first) */
	
	/* Open music file in new emulator */
	Music_Emu* emu;
	if( handle_error( gme_open_file( src, &emu, sample_rate ) ) )
  {
    return 1;
  }
	
	/* Start track */
	if( handle_error( gme_start_track( emu, track ) ) ) 
  {
    return 2;
  }
	
	/* Begin writing to wave file */
	wave_open( sample_rate, dest );
	wave_enable_stereo();
	
	/* Record 10 seconds of track */
	while ( gme_track_ended( emu ) == 0 )
	{
		/* Sample buffer */
		#define buf_size 1024 /* can be any multiple of 2 */
		short buf [buf_size];
		
		/* Fill sample buffer */
		if ( handle_error( gme_play( emu, buf_size, buf ) ) )
    {
      return 3;
    }
		
		/* Write samples to wave file */
		wave_write( buf, buf_size );
	}
	
	/* Cleanup */
	gme_delete( emu );
	wave_close();
	
	return 0;
}

int handle_error( const char* str )
{
	if ( str )
	{
		printf( "Error: %s\n", str ); getchar();
		return 1;
	}
  return 0;
}
